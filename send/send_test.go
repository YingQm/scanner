package send

import (
	"github.com/YingQm/scanner/config"
	"testing"
)

func TestSendEmail(t *testing.T) {
	var cfg config.Config

	cfg.FromEmail = "fuzamei@126.com"
	cfg.FromEmailPsw = "授权码"
	cfg.ToEmail = "yqm@disanbo.com,2355241170@qq.com"
	cfg.Host = "smtp.126.com"
	cfg.PostEmail = 25

	SendEmail(&cfg, "测试", "测试发送")
}
