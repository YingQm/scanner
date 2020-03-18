package send

import (
<<<<<<< HEAD
	"gitlab.33.cn/wallet/monitor/types"
=======
	"github.com/YingQm/scanner/config"
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
	"testing"
)

func TestSendEmail(t *testing.T) {
<<<<<<< HEAD
	var cfg types.Config

	cfg.FromEmail = "fuzamei@126.com"
	cfg.FromEmailPsw = "fzm123456"
=======
	var cfg config.Config

	cfg.FromEmail = "fuzamei@126.com"
	cfg.FromEmailPsw = "授权码"
>>>>>>> 43b7f57b001bf47c96bca9ea10c9a48e391176fb
	cfg.ToEmail = "yqm@disanbo.com,2355241170@qq.com"
	cfg.Host = "smtp.126.com"
	cfg.PostEmail = 25

	SendEmail(&cfg, "测试", "测试发送")
}
