package main

import (
	"log"
	"net/http"

	"github.com/vladkonst/metrics-alerting/internal/configs"
	"github.com/vladkonst/metrics-alerting/internal/models"
	"github.com/vladkonst/metrics-alerting/internal/storage"
	"github.com/vladkonst/metrics-alerting/routers"
)

func main() {
	cfg := configs.GetServerConfig()
	metricsCh := make(chan models.Metrics)
	fileStorage, err := storage.NewFileManager(cfg.IntervalsCfg.FileStoragePath, cfg.IntervalsCfg.Restore, cfg.IntervalsCfg.StoreInterval, &metricsCh)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		if err := fileStorage.LoadMetrics(); err != nil {
			log.Panic(err)
		}
	}()

	http.ListenAndServe(cfg.NetAddressCfg.String(), routers.GetRouter())
}
