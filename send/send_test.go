package send

import (
	"gitlab.33.cn/wallet/monitor/types"
	"testing"
)

func TestSendEmail(t *testing.T) {
	var cfg types.Config

	cfg.FromEmail = "fuzamei@126.com"
	cfg.FromEmailPsw = "fzm123456"
	cfg.ToEmail = "yqm@disanbo.com,2355241170@qq.com"
	cfg.Host = "smtp.126.com"
	cfg.PostEmail = 25

	SendEmail(&cfg, "测试", "测试发送")
}
