package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type Metrics struct {
	mu        sync.Mutex
	Gauges    map[string]float64
	PollCount int
}

func (m *Metrics) updateGaugeMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.mu.Lock()
	m.Gauges = map[string]float64{
		"Alloc":         float64(memStats.Alloc),
		"BuckHashSys":   float64(memStats.BuckHashSys),
		"Frees":         float64(memStats.Frees),
		"GCCPUFraction": float64(memStats.GCCPUFraction),
		"GCSys":         float64(memStats.GCSys),
		"HeapAlloc":     float64(memStats.HeapAlloc),
		"HeapIdle":      float64(memStats.HeapIdle),
		"HeapInuse":     float64(memStats.HeapInuse),
		"HeapObjects":   float64(memStats.HeapObjects),
		"HeapReleased":  float64(memStats.HeapReleased),
		"HeapSys":       float64(memStats.HeapSys),
		"LastGC":        float64(memStats.LastGC),
		"Lookups":       float64(memStats.Lookups),
		"MCacheInuse":   float64(memStats.MCacheInuse),
		"MCacheSys":     float64(memStats.MCacheSys),
		"MSpanInuse":    float64(memStats.MSpanInuse),
		"MSpanSys":      float64(memStats.MSpanSys),
		"Mallocs":       float64(memStats.Mallocs),
		"NextGC":        float64(memStats.NextGC),
		"NumForcedGC":   float64(memStats.NumForcedGC),
		"NumGC":         float64(memStats.NumGC),
		"OtherSys":      float64(memStats.OtherSys),
		"PauseTotalNs":  float64(memStats.PauseTotalNs),
		"StackInuse":    float64(memStats.StackInuse),
		"StackSys":      float64(memStats.StackSys),
		"Sys":           float64(memStats.Sys),
		"TotalAlloc":    float64(memStats.TotalAlloc),
		"RandomValue":   1.0 + rand.Float64()*9,
	}
	m.PollCount++
	m.mu.Unlock()
}

func sendGaugeMetrics(m *Metrics) {
	for k, v := range m.Gauges {
		resp, err := http.Post(fmt.Sprintf("http://localhost:8080/update/gauge/%v/%v", k, v), "text/plain", nil)
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

func sendCounterMetrics(m *Metrics) {
	resp, err := http.Post(fmt.Sprintf("http://localhost:8080/update/counter/PollCount/%v", m.PollCount), "text/plain", nil)
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
	pollInterval := time.Duration(2)
	reportInterval := time.Duration(4)
	metrics := Metrics{Gauges: make(map[string]float64)}
	ticker := time.NewTicker(reportInterval * time.Second)
	go func() {
		ticker := time.NewTicker(pollInterval * time.Second)
		for range ticker.C {
			metrics.updateGaugeMetrics()
		}
	}()

	for range ticker.C {
		sendGaugeMetrics(&metrics)
		sendCounterMetrics(&metrics)
	}

}
