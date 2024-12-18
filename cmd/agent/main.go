package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vladkonst/metrics-alerting/internal/agent"
	"github.com/vladkonst/metrics-alerting/internal/configs"
	"github.com/vladkonst/metrics-alerting/internal/models"
)

var timings = []time.Duration{0, time.Second, time.Second * 3, time.Second * 5}

func sendRequest(m map[string]models.Metrics, serverAddr *configs.NetAddressCfg, tryCount int) {
	if tryCount == 4 {
		return
	}

	metrics := make([]models.Metrics, 0)
	for _, v := range m {
		metrics = append(metrics, v)
	}

	b, err := json.Marshal(metrics)

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

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/updates/", serverAddr.String()), buff)
	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	if tryCount > 0 {
		time.Sleep(timings[tryCount])
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		var opError *net.OpError
		if errors.As(err, &opError) && opError.Op == "dial" {
			go sendRequest(m, serverAddr, tryCount+1)
		}
		return
	}

	resBody, _ := io.ReadAll(resp.Body)
	log.Println("Body: ", string(resBody))
	log.Println("Status: ", resp.StatusCode)
	log.Println("Content-Type: ", resp.Header.Get("Content-Type"))
	log.Println("Content-Encoding: ", resp.Header.Values("Content-Encoding"))
	resp.Body.Close()
}

func sendMetrics(cfg *configs.ClientCfg, metricsCh *chan models.Metrics) {
	reprotTicker := time.NewTicker(time.Duration(cfg.IntervalsCfg.ReportInterval) * time.Second)
	metrics := make(map[string]models.Metrics)
	for {
		select {
		case <-reprotTicker.C:
			sendRequest(metrics, cfg.NetAddressCfg, 0)
		case metric := <-*metricsCh:
			metrics[metric.ID] = metric
		}
	}
}

func main() {
	done := make(chan bool)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		done <- true
	}()
	cfg := configs.GetClientConfig()
	metricsStorage := agent.NewMetricsStorage()
	metricsStorage.InitMetrics()
	pollTicker := time.NewTicker(time.Duration(cfg.IntervalsCfg.PollInterval) * time.Second)
	metricsCh := make(chan models.Metrics)

	go func() {
		sendMetrics(cfg, &metricsCh)
	}()

	for {
		select {
		case <-pollTicker.C:
			metricsStorage.UpdateMetrics(&metricsCh)
		case <-done:
			return
		}
	}
}
