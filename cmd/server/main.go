package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vladkonst/metrics-alerting/internal/configs"
	"github.com/vladkonst/metrics-alerting/internal/storage"
	"github.com/vladkonst/metrics-alerting/routers"
)

func main() {
	done := make(chan bool)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		done <- true
	}()

	cfg := configs.GetServerConfig()
	memStorage := storage.GetStorage()
	fileStorage, err := storage.NewFileManager(cfg.IntervalsCfg.FileStoragePath, cfg.IntervalsCfg.Restore, cfg.IntervalsCfg.StoreInterval, memStorage.GetMetricsChanel())
	if err != nil {
		log.Panic(err)
	}

	go func() {
		if err := fileStorage.ProcessMetrics(); err != nil {
			log.Panic(err)
		}
	}()

	go func() {
		log.Panic(http.ListenAndServe(cfg.NetAddressCfg.String(), routers.GetRouter(memStorage)))
	}()

	<-done
	fileStorage.LoadMetrics()
}
