package util

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	l "github.com/inconshreveable/log15"
	"gitlab.33.cn/wallet/monitor/types"
	"os"
)

var log = l.New("module", "util")

type DbHandler struct {
	Db  *sql.DB
	Cfg *types.Config
}

func NewMysql(cfg *types.Config) (*DbHandler, error) {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", "root", cfg.Mysql.MysqlPwd,
		cfg.Mysql.MysqlIp, cfg.Mysql.MysqlPort, cfg.Mysql.DbName)
	log.Info("NewMysql", "url:", url)

	db, err := sql.Open("mysql", url)
	if err != nil {
		log.Error("mysql", "open error", err.Error())
		os.Exit(0)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)

	err = db.Ping()
	if err != nil {
		log.Error("mysql", "ping error", err.Error())
		return nil, err
	}

	if cfg.Mysql.NodePortTable != "" {
		sqlstr := fmt.Sprintf(nodeportinfo, cfg.Mysql.NodePortTable)
		_, err = db.Exec(sqlstr)
		if err != nil {
			return nil, err
		}

	}
	if cfg.Mysql.NodeSyncTable != "" {
		sqlstr := fmt.Sprintf(nodesyncinfo, cfg.Mysql.NodeSyncTable)
		_, err = db.Exec(sqlstr)
		if err != nil {
			return nil, err
		}
	}

	xdb := new(DbHandler)
	xdb.Db = db
	xdb.Cfg = cfg
	return xdb, nil
}

//获取地址
func (m *DbHandler) FetchNodePortAddr() map[string]string {

	cmdstr := fmt.Sprintf(`select cointype, nodeaddr  from %s`, m.Cfg.Mysql.NodePortTable)
	rows, err := m.Db.Query(cmdstr)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	defer rows.Close()
	addrMap := make(map[string]string)

	var cointype, addrs string
	for rows.Next() {
		err := rows.Scan(&cointype, &addrs)
		if err != nil {
			log.Error("Scan", "err:%v", err.Error())
			continue
		}

		addrMap[cointype] = addrs
	}
	return addrMap

}

//获取地址
func (m *DbHandler) FetchNodeSyncAddr() []types.NodeSync {

	cmdstr := fmt.Sprintf(`select cointype, nodeaddr,targetaddr,nodetype,maxdiff  from %s`, m.Cfg.Mysql.NodeSyncTable)
	rows, err := m.Db.Query(cmdstr)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	defer rows.Close()
	nodeAddrs := []types.NodeSync{}

	var node types.NodeSync
	for rows.Next() {
		err := rows.Scan(&node.CoinType, &node.Nodes, &node.Target, &node.NodeType, &node.MaxDiff)
		if err != nil {
			log.Error("Scan", "err:%v", err.Error())
			continue
		}

		nodeAddrs = append(nodeAddrs, node)
	}
	return nodeAddrs

}
