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
var thread chan int = make(chan int)
var nowThread int

func RunScan(ipports []string) {
	for {
		scanResult = ""
		nowThread = len(ipports)
		fmt.Println("len(ipports)", nowThread)
		for i := 0; i < len(ipports); i++ {
			go scan(ipports[i])
		}

		<-thread
		fmt.Println("一次循环结束", scanResult)
		if len(scanResult) > 0 && !StringSliceReflectEqual(scanResult, scanResultOld) {
			sendEmail("提醒", "端口不开放:"+scanResult)
		}
		scanResultOld = scanResult

		time.Sleep((time.Duration)(cfg.IntervalTime) * time.Minute)
	}
}

func StringSliceReflectEqual(a, b string) bool {
	//	return reflect.DeepEqual(strings.Split(a, ","), strings.Split(b, ","))
	return len(a) == len(b)
}

func scan(address string) {
	_, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(address)
		result <- address
	}

	nowThread--
	fmt.Println("nowThread", nowThread)
	if nowThread == 0 {
		time.Sleep(time.Second)
		thread <- 0
	}
}

func getResult() {
	for {
		s, ok := <-result
		if ok {
			scanResult += (string)(s + ";")
		} else {
			fmt.Println("err: ", ok)
		}
	}
	fmt.Println("/r/n err 退出了 getResult /r/n")
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
	go getResult()

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
	fmt.Println("发送邮件", body)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("邮件发送失败", err)
	} else {
		fmt.Println("邮件发送成功")
	}
}
