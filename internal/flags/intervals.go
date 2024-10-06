package flags

import (
	"github.com/caarlos0/env"

	"flag"
	"fmt"
)

type IntervalsCfg struct {
	ReportInterval int `env:"REPORT_INTERVAL"`
	PollInterval   int `env:"POLL_INTERVAL"`
}

func GetIntervalsCfg() *IntervalsCfg {
	intervalCfg := &IntervalsCfg{ReportInterval: 10, PollInterval: 2}
	flag.IntVar(&intervalCfg.ReportInterval, "r", intervalCfg.ReportInterval, "report interval to send metrics")
	flag.IntVar(&intervalCfg.PollInterval, "p", intervalCfg.PollInterval, "poll interval to update metrics")
	flag.Parse()
	if err := env.Parse(intervalCfg); err != nil {
		fmt.Println("can't parse intervals from env variables")
	}
	return intervalCfg
}
