package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vladkonst/metrics-alerting/internal/configs"
	"github.com/vladkonst/metrics-alerting/internal/metrics"
)

func sendGaugeMetrics(serverAddr *configs.NetAddressCfg, m *metrics.Metrics) {
	for k, v := range m.Gauges {
		resp, err := http.Post(fmt.Sprintf("http://%s/update/gauge/%v/%v", serverAddr.String(), k, v), "text/plain", nil)
		if err != nil {
			log.Fatal(err)
		}
		resBody, _ := io.ReadAll(resp.Body)
		log.Println(string(resBody))
		log.Println(resp.StatusCode)
		log.Println(resp.Header.Get("Content-Type"))
		resp.Body.Close()
	}
}

func sendCounterMetrics(serverAddr *configs.NetAddressCfg, m *metrics.Metrics) {
	resp, err := http.Post(fmt.Sprintf("http://%s/update/counter/PollCount/%v", serverAddr.String(), m.PollCount), "text/plain", nil)
	if err != nil {
		log.Fatal(err)
	}
	resBody, _ := io.ReadAll(resp.Body)
	log.Println(string(resBody))
	log.Println(resp.StatusCode)
	log.Println(resp.Header.Get("Content-Type"))
	resp.Body.Close()
}

func main() {
	cfg := configs.GetClientConfig()
	metrics := metrics.Metrics{Gauges: make(map[string]float64)}
	ticker := time.NewTicker(time.Duration(cfg.IntervalsCfg.ReportInterval) * time.Second)
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.IntervalsCfg.PollInterval) * time.Second)
		for range ticker.C {
			metrics.UpdateGaugeMetrics()
		}
	}()

	for range ticker.C {
		sendGaugeMetrics(cfg.NetAddressCfg, &metrics)
		sendCounterMetrics(cfg.NetAddressCfg, &metrics)
	}
}
