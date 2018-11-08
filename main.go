package main

import (
	"flag"
	"fmt"
	"github.com/go-gomail/gomail"
	"mytest/scanner/config"
	"net"
	"strings"
	"time"
)

//关闭程序
var clo chan bool = make(chan bool)
var cfg *config.Config
var result chan string = make(chan string)
var scanResultOld string
var scanResult string

func RunScan(ipports []string) {
	for {
		scanResult = ""
		for i := 0; i < len(ipports); i++ {
			go scan(ipports[i])
		}

		for i := 0; i < len(ipports); i++ {
			s, ok := <-result
			if ok && len(s) > 0 {
				scanResult += (string)(s + "\n")
			}
		}

		fmt.Println("一次循环结束")
		if len(scanResult) > 0 && !StringSliceReflectEqual(scanResult, scanResultOld) {
			sendEmail("提醒", "端口不开放:"+scanResult)
		}
		scanResultOld = scanResult

		time.Sleep((time.Duration)(cfg.IntervalTime) * time.Minute)
	}
}

func StringSliceReflectEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	as := strings.Split(a, "\n")
	bs := strings.Split(b, "\n")
	if len(as) != len(bs) {
		return false
	}

	for i := 0; i < len(as); i++ {
		bfind := false
		for j := 0; j < len(bs); j++ {
			if as[i] == bs[j] {
				bfind = true
				bs = append(bs[:j], bs[j+1:]...)
				break
			}
		}

		if !bfind {
			return false
		}
	}

	return true
}

func scan(address string) {
	_, err := net.Dial("tcp", address)
	if err != nil {
		//	fmt.Println(address)
		result <- address
	} else {
		result <- ""
	}
}

func SendToEmail() {
	for {
		if len(scanResultOld) > 0 {
			fmt.Println("scanResultOld", scanResultOld)
			sendEmail("提醒", "端口不开放:"+scanResultOld)
		}
		time.Sleep((time.Duration)(cfg.SendTime) * time.Minute)
	}
}

func main() {
	configpath := flag.String("f", "config.toml", "configfile")
	cfg = config.InitCfg(*configpath)

	go RunScan(cfg.IpPosts)
	go SendToEmail()
	//等待退出指令
	<-clo
	fmt.Println("Exit")
}

func sendEmail(title, body string) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", cfg.FromEmail, cfg.FromEmail) // 发件人

	// 收件人
	var toEmails []string = make([]string, 0)
	tos := strings.Split(cfg.ToEmail, ",")
	for i := 0; i < len(tos); i++ {
		to := m.FormatAddress(tos[i], tos[i])
		toEmails = append(toEmails, to)
	}
	toMap := make(map[string][]string)
	toMap["To"] = toEmails
	m.SetHeaders(toMap)

	m.SetHeader("Subject", title) // 主题
	m.SetBody("text/html", body)  // 正文

	d := gomail.NewPlainDialer(cfg.Host, (int)(cfg.PostEmail), cfg.FromEmail, cfg.FromEmailPsw) // 发送邮件服务器、端口、发件人账号、发件人密码
	fmt.Println("发送邮件")
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("邮件发送失败", err)
	} else {
		fmt.Println("邮件发送成功")
	}
}
