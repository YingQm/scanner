package config

import (
	"fmt"
	tml "github.com/BurntSushi/toml"
	"net"
	"os"
	"strconv"
	"strings"
)

func InitCfg(path string) *Config {
	var cfg Config
	if _, err := tml.DecodeFile(path, &cfg); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	ips := processIp(cfg.IpAddrs)
	ports := processPort(cfg.Ports)
	ipports := make([]string, 0)
	for i := 0; i < len(ips); i++ {
		for j := 0; j < len(ports); j++ {
			ipports = append(ipports, ips[i]+":"+strconv.Itoa(ports[j]))
		}
	}
	cfg.IpPosts = ipports

	if cfg.IntervalTime < 1 {
		cfg.IntervalTime = 1
	}

	if cfg.SendTime < 10 {
		cfg.SendTime = 10
	}

	fromEmail := strings.ToLower(cfg.FromEmail)
	if strings.Contains(fromEmail, "@qq.com") {
		cfg.Host = "smtp.qq.com"
		cfg.PostEmail = 465
	} else if strings.Contains(fromEmail, "@126.com") {
		cfg.Host = "smtp.126.com"
		cfg.PostEmail = 25
	} else if strings.Contains(fromEmail, "@163.com") {
		cfg.Host = "smtp.163.com"
		cfg.PostEmail = 25
	}

	return &cfg
}

func processIp(IpAddrs string) []string {
	var ips = make([]string, 0)
	ipso := strings.Split(IpAddrs, ",")
	for i := 0; i < len(ipso); i++ {
		si := net.ParseIP(ipso[i])
		if si == nil {
			continue
		}
		ips = append(ips, ipso[i])
	}
	return ips
}

func processPort(Ports string) []int {
	var ports []int = make([]int, 0)
	ps := strings.Split(Ports, ",")
	for i := 0; i < len(ps); i++ {
		p, err := strconv.Atoi(ps[i])
		if err != nil {
			continue
		}

		if p >= 1 && p <= 65535 {
			ports = append(ports, p)
		}
	}
	return ports
}