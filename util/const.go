package util

const (
	nodeportinfo = `create table if not exists %s (
		cointype varchar(64) not null comment '币种',
		nodeaddr varchar(512) not null comment '格式ip:port,多个地址以逗号,分隔',
		primary key (cointype))  default character set utf8;`
	nodesyncinfo = `create table if not exists %s (
		cointype varchar(64) not null comment '币种',
		nodeaddr varchar(512) not null comment '格式ip:port,多个地址以逗号,分隔',
		targetaddr varchar (64) not null comment '主网地址',
		nodetype  tinyint(4) not null default '0' comment '节点类型 0-普通节点 1-BTC insight-api',
		maxdiff   int(11) not null default '0' comment '与目标节点相比允许的最大高度差',
		primary key (cointype,nodetype))  default character set utf8;`
)
