package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

func getGaugeMetrics() map[string]float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	gauges := map[string]float64{
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
	return gauges
}

func main() {
	PollCount := 0
	pollInterval := 2
	// reportInterval := 10
	for {
		time.Sleep(time.Second * time.Duration(pollInterval))
		gaugeMetrics := getGaugeMetrics()
		PollCount++
		if PollCount%5 == 0 {
			for k, v := range gaugeMetrics {
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
			resp, err := http.Post(fmt.Sprintf("http://localhost:8080/update/counter/PollCount/%v", PollCount), "text/plain", nil)
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

}
