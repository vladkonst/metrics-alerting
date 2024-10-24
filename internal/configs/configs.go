package configs

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

type ClientCfg struct {
	IntervalsCfg  *IntervalsCfg
	NetAddressCfg *NetAddressCfg
}

func GetClientConfig() *ClientCfg {
	intervalCfg := &IntervalsCfg{ReportInterval: 10, PollInterval: 2}
	flag.IntVar(&intervalCfg.ReportInterval, "r", intervalCfg.ReportInterval, "report interval to send metrics")
	flag.IntVar(&intervalCfg.PollInterval, "p", intervalCfg.PollInterval, "poll interval to update metrics")
	addr := &NetAddressCfg{Host: "localhost", Port: 8080}
	flag.Var(addr, "a", "Server net address host:port")
	flag.Parse()
	if err := env.Parse(intervalCfg); err != nil {
		fmt.Println("can't parse intervals from env variables")
	}
	if adr := os.Getenv("ADDRESS"); adr != "" {
		addr.Set(os.Getenv("ADDRESS"))
	}
	return &ClientCfg{IntervalsCfg: intervalCfg, NetAddressCfg: addr}
}

func GetServerConfig() *NetAddressCfg {
	addr := &NetAddressCfg{Host: "localhost", Port: 8080}
	flag.Var(addr, "a", "Server net address host:port")
	flag.Parse()
	if adr := os.Getenv("ADDRESS"); adr != "" {
		addr.Set(os.Getenv("ADDRESS"))
	}
	return addr
}
