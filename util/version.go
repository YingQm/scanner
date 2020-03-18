package util

const version = "1.0.0"

var (
	GitCommit string
	BuildTime string
)

//GetVersion 获取版本信息
func GetVersion() string {
	if GitCommit != "" {
		return version + "-" + GitCommit + "-" + BuildTime
	}
	return version
}
