package send

import (
	"errors"
<<<<<<< HEAD
	"github.com/go-gomail/gomail"
	l "github.com/inconshreveable/log15"
	"gitlab.33.cn/wallet/monitor/types"
	"strings"
)

var log = l.New("module", "send")

func SendEmail(cfg *types.Config, title, body string) error {
=======
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/YingQm/scanner/config"
	"strings"
)

func SendEmail(cfg *config.Config, title, body string) error {
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
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
<<<<<<< HEAD
	log.Info("发送邮件")
	log.Info("SendEmail", "content", body)
	if err := d.DialAndSend(m); err != nil {
		log.Info("SendEmail", "邮件发送失败", err.Error())
		return err
	} else {
		log.Info("邮件发送成功")
=======
	fmt.Println("发送邮件", body)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("邮件发送失败", err)
		return err
	} else {
		fmt.Println("邮件发送成功")
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
	}

	return nil
}
