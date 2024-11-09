package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vladkonst/metrics-alerting/internal/agent"
	"github.com/vladkonst/metrics-alerting/internal/configs"
	"github.com/vladkonst/metrics-alerting/internal/models"
)

func sendRequest(v *models.Metrics, serverAddr *configs.NetAddressCfg) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		return
	}

	buff := bytes.NewBuffer(nil)
	zb := gzip.NewWriter(buff)
	_, err = zb.Write(b)
	if err != nil {
		log.Println(err)
		return
	}

	err = zb.Close()
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/update/", serverAddr.String()), buff)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	resBody, _ := io.ReadAll(resp.Body)
	log.Println("Body: ", string(resBody))
	log.Println("Status: ", resp.StatusCode)
	log.Println("Content-Type: ", resp.Header.Get("Content-Type"))
	log.Println("Content-Encoding: ", resp.Header.Values("Content-Encoding"))
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
