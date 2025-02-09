package configs

import (
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env"
)

type ClientCfg struct {
	IntervalsCfg  *ClientIntervalsCfg
	NetAddressCfg *NetAddressCfg
}

type ServerCfg struct {
	IntervalsCfg  *ServerIntervalsCfg
	NetAddressCfg *NetAddressCfg
}

func GetClientConfig() *ClientCfg {
	intervalCfg := &ClientIntervalsCfg{ReportInterval: 3, PollInterval: 1}
	flag.IntVar(&intervalCfg.ReportInterval, "r", intervalCfg.ReportInterval, "report interval to send metrics")
	flag.IntVar(&intervalCfg.PollInterval, "p", intervalCfg.PollInterval, "poll interval to update metrics")
	flag.IntVar(&intervalCfg.RateLimit, "l", 1, "requests rate limit number")
	flag.StringVar(&intervalCfg.HashKey, "k", "", "hash key")
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

func GetServerConfig() *ServerCfg {
	addr := &NetAddressCfg{Host: "localhost", Port: 8080}
	intervalCfg := &ServerIntervalsCfg{StoreInterval: 300, FileStoragePath: "metrics.txt", Restore: true}
	flag.Var(addr, "a", "Server net address host:port")
	flag.IntVar(&intervalCfg.StoreInterval, "i", intervalCfg.StoreInterval, "store interval to load metrics to the file")
	flag.StringVar(&intervalCfg.FileStoragePath, "f", intervalCfg.FileStoragePath, "file with stored metrics")
	flag.StringVar(&intervalCfg.DatabaseDSN, "d", "", "database connection string")
	flag.StringVar(&intervalCfg.HashKey, "k", "", "hash key")
	flag.BoolVar(&intervalCfg.Restore, "r", intervalCfg.Restore, "allow metrics load from file on server start")
	flag.Parse()
	if err := env.Parse(intervalCfg); err != nil {
		fmt.Println("can't parse intervals from env variables")
	}

	if adr := os.Getenv("ADDRESS"); adr != "" {
		addr.Set(os.Getenv("ADDRESS"))
	}

	return &ServerCfg{IntervalsCfg: intervalCfg, NetAddressCfg: addr}
}
