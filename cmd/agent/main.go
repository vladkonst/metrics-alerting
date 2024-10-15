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
	for k, v := range m.Counters {
		resp, err := http.Post(fmt.Sprintf("http://%s/update/counter/%v/%v", serverAddr.String(), k, v), "text/plain", nil)
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

func main() {
	cfg := configs.GetClientConfig()
	metrics := metrics.Metrics{Gauges: make(map[string]float64), Counters: make(map[string]int)}
	reprotTicker := time.NewTicker(time.Duration(cfg.IntervalsCfg.ReportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(cfg.IntervalsCfg.PollInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			metrics.UpdateGaugeMetrics()
		default:
		}
		select {
		case <-reprotTicker.C:
			sendGaugeMetrics(cfg.NetAddressCfg, &metrics)
			sendCounterMetrics(cfg.NetAddressCfg, &metrics)
		default:
		}
	}
}
