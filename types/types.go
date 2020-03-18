package types

type Config struct {
	Title        string `toml:"title"`
	Log          *Log
	Mysql        mysql
	IntervalTime int64  `toml:"intervalTime"`
	FromEmail    string `toml:"fromEmail"`
	FromEmailPsw string `toml:"fromEmailPsw"`
	ToEmail      string `toml:"toEmail"`
	Host         string `toml:"host"`
	PostEmail    int64  `toml:"postEmail"`
	SendTime     int64  `toml:"sendTime"`
	Rpcname      string `toml:"noderpcuser"`
	Rpcpasswd    string `toml:"noderpcpasswd"`
	Omniname     string `toml:"omniuser"`
	Omnipasswd   string `toml:"omnipasswd"`
	RateLimit    RateLimit
	Parallel     string `toml:"parallel"`
}

type mysql struct {
	MysqlIp       string `toml:"mysqlIp"`
	MysqlPort     int    `toml:"mysqlPort"`
	MysqlPwd      string `toml:"mysqlPwd"`
	DbName        string `toml:"dbName"`
	NodePortTable string `toml:"nodePortTable"`
	NodeSyncTable string `toml:"nodeSyncTable"`
}

type Message struct {
	CoinType string
	Addr     string
	ErrMsg   string
}

type NodePort struct {
	CoinType string
	Nodes    string
}

type NodeSync struct {
	CoinType string
	Nodes    string
	Target   string
	NodeType int64
	MaxDiff  int64
}

type ClientResponse struct {
	Id     uint64      `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

//限流配置
type RateLimit struct {
	LimitOpen        bool  `toml:"limitOpen"`
	MaxConCurrentNum int   `json:"maxConCurrentNum"`
	TimeInterval     int64 `toml:"timeInterval"`
	MaxCount         int64 `toml:"maxCount"`
}
