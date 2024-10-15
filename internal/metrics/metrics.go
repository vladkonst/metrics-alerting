package metrics

import (
	"math/rand/v2"
	"runtime"
	"sync"
)

type Metrics struct {
	mu       sync.Mutex
	Gauges   map[string]float64
	Counters map[string]int
}

func (m *Metrics) UpdateGaugeMetrics() {
	m.mu.Lock()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	m.Gauges["Alloc"] = float64(memStats.Alloc)
	m.Gauges["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.Gauges["Frees"] = float64(memStats.Frees)
	m.Gauges["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	m.Gauges["GCSys"] = float64(memStats.GCSys)
	m.Gauges["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.Gauges["HeapIdle"] = float64(memStats.HeapIdle)
	m.Gauges["HeapInuse"] = float64(memStats.HeapInuse)
	m.Gauges["HeapObjects"] = float64(memStats.HeapObjects)
	m.Gauges["HeapReleased"] = float64(memStats.HeapReleased)
	m.Gauges["HeapSys"] = float64(memStats.HeapSys)
	m.Gauges["LastGC"] = float64(memStats.LastGC)
	m.Gauges["Lookups"] = float64(memStats.Lookups)
	m.Gauges["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.Gauges["MCacheSys"] = float64(memStats.MCacheSys)
	m.Gauges["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.Gauges["MSpanSys"] = float64(memStats.MSpanSys)
	m.Gauges["Mallocs"] = float64(memStats.Mallocs)
	m.Gauges["NextGC"] = float64(memStats.NextGC)
	m.Gauges["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.Gauges["NumGC"] = float64(memStats.NumGC)
	m.Gauges["OtherSys"] = float64(memStats.OtherSys)
	m.Gauges["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.Gauges["StackInuse"] = float64(memStats.StackInuse)
	m.Gauges["StackSys"] = float64(memStats.StackSys)
	m.Gauges["Sys"] = float64(memStats.Sys)
	m.Gauges["TotalAlloc"] = float64(memStats.TotalAlloc)
	m.Gauges["RandomValue"] = 1.0 + rand.Float64()*9
	m.Counters["PollCount"]++
	m.mu.Unlock()
}
