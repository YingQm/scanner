package send

import (
	"errors"
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/YingQm/scanner/config"
	"strings"
)

func SendEmail(cfg *config.Config, title, body string) error {
	if cfg == nil {
		return errors.New("配置文件错误，配置为空！")
	}

	if len(cfg.FromEmail) == 0 || len(cfg.FromEmailPsw) == 0 || len(cfg.ToEmail) == 0 || len(cfg.Host) == 0 {
		return errors.New("配置文件错误，没有配置发送邮件信息！")
	}

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
		return err
	} else {
		fmt.Println("邮件发送成功")
	}

	return nil
}
