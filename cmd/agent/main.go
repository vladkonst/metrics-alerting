package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vladkonst/metrics-alerting/internal/agent"
	"github.com/vladkonst/metrics-alerting/internal/configs"
)

func sendRequest(v interface{}, serverAddr *configs.NetAddressCfg) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
	}

	buff := bytes.NewBuffer(b)
	resp, err := http.Post(fmt.Sprintf("http://%s/update/", serverAddr.String()), "encoding/json", buff)
	if err != nil {
		log.Fatal(err)
	}
	resBody, _ := io.ReadAll(resp.Body)
	log.Println(string(resBody))
	log.Println(resp.StatusCode)
	log.Println(resp.Header.Get("Content-Type"))
	resp.Body.Close()
}

func sendMetrics(serverAddr *configs.NetAddressCfg, m *agent.MetricsStorage) {
	for _, v := range m.Gauges {
		sendRequest(v, serverAddr)
	}

	for _, v := range m.Counters {
		sendRequest(v, serverAddr)
	}
}

func main() {
	cfg := configs.GetClientConfig()
	metricsStorage := agent.NewMetricsStorage()
	metricsStorage.InitMetrics()
	reprotTicker := time.NewTicker(time.Duration(cfg.IntervalsCfg.ReportInterval) * time.Second)
	pollTicker := time.NewTicker(time.Duration(cfg.IntervalsCfg.PollInterval) * time.Second)

	for {
		select {
		case <-pollTicker.C:
			metricsStorage.UpdateMetrics()
		default:
		}
		select {
		case <-reprotTicker.C:
			sendMetrics(cfg.NetAddressCfg, &metricsStorage)
		default:
		}
	}
}
