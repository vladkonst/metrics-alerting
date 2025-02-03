package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
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

type hasher struct {
	hash hash.Hash
}

func NewHasher(key string) *hasher {
	if key == "" {
		return nil
	}
	hash := sha256.New()
	return &hasher{hash}
}

func (h *hasher) hashBody(body []byte) []byte {
	_, err := h.hash.Write(body)
	if err != nil {
		log.Panic("failed to hash request body")
	}

	dst := h.hash.Sum(nil)
	return dst
}

func sendRequest(metricsJobs chan models.Metrics, serverAddr *configs.NetAddressCfg, tryCount int, h *hasher) {
	if tryCount == 4 {
		return
	}

	metrics := make([]models.Metrics, 0)
	for m := range metricsJobs {
		metrics = append(metrics, m)
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

	// if h != nil {
	// 	hashedBody := h.hashBody(b)
	// 	req.Header.Set("HashSHA256", hex.EncodeToString(hashedBody))
	// }

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
			metricsJobs := make(chan models.Metrics)
			for _, m := range metrics {
				metricsJobs <- m
			}
			close(metricsJobs)
			go sendRequest(metricsJobs, serverAddr, tryCount+1, h)
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

func sendMetrics(cfg *configs.ClientCfg, metricsCh *chan models.Metrics, done chan struct{}, h *hasher) {
	reprotTicker := time.NewTicker(time.Duration(cfg.IntervalsCfg.ReportInterval) * time.Second)
	metrics := make([]models.Metrics, 0)
	for {
		select {
		case <-done:
			return
		case <-reprotTicker.C:
			metricsJobs := make(chan models.Metrics, len(metrics))
			for _, metric := range metrics {
				metricsJobs <- metric
			}
			close(metricsJobs)
			for i := 0; i < cfg.IntervalsCfg.RateLimit; i++ {
				sendRequest(metricsJobs, cfg.NetAddressCfg, 0, h)
			}
		case metric := <-*metricsCh:
			metrics = append(metrics, metric)
		}
	}
}

func main() {
	done := make(chan struct{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		close(done)
	}()
	cfg := configs.GetClientConfig()
	metricsStorage := agent.NewMetricsStorage()
	metricsStorage.InitMetrics()
	metricsCh := make(chan models.Metrics)
	h := NewHasher(cfg.IntervalsCfg.HashKey)

	go sendMetrics(cfg, &metricsCh, done, h)

	go func(done chan struct{}) {
		pollTicker := time.NewTicker(time.Duration(cfg.IntervalsCfg.PollInterval) * time.Second)
		for {
			select {
			case <-done:
				return
			case <-pollTicker.C:
				metricsStorage.UpdateRuntimeMetrics(&metricsCh)
			}
		}
	}(done)

	go func(done chan struct{}) {
		pollTicker := time.NewTicker(time.Duration(cfg.IntervalsCfg.PollInterval) * time.Second)
		for {
			select {
			case <-done:
				return
			case <-pollTicker.C:
				err := metricsStorage.UpdatePSUtilMetrics(&metricsCh)
				if err != nil {
					log.Println(err)
					close(done)
				}
			}
		}
	}(done)

	<-done
	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	default:
	// 	}
	// }
}
