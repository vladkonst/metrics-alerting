package agent

import (
	"math/rand/v2"
	"runtime"
	"sync"

	"github.com/vladkonst/metrics-alerting/internal/models"
)

type MetricsStorage struct {
	mu       sync.Mutex
	once     sync.Once
	Gauges   map[string]*models.Metrics
	Counters map[string]*models.Metrics
}

func NewMetricsStorage() (ms MetricsStorage) {
	ms = MetricsStorage{Gauges: make(map[string]*models.Metrics), Counters: make(map[string]*models.Metrics)}
	return
}

func (m *MetricsStorage) UpdateMetrics() {
	m.mu.Lock()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	Alloc := float64(memStats.Alloc)
	BuckHashSys := float64(memStats.BuckHashSys)
	Frees := float64(memStats.Frees)
	GCCPUFraction := float64(memStats.GCCPUFraction)
	GCSys := float64(memStats.GCSys)
	HeapAlloc := float64(memStats.HeapAlloc)
	HeapIdle := float64(memStats.HeapIdle)
	HeapInuse := float64(memStats.HeapInuse)
	HeapObjects := float64(memStats.HeapObjects)
	HeapReleased := float64(memStats.HeapReleased)
	HeapSys := float64(memStats.HeapSys)
	LastGC := float64(memStats.LastGC)
	Lookups := float64(memStats.Lookups)
	MCacheInuse := float64(memStats.MCacheInuse)
	MCacheSys := float64(memStats.MCacheSys)
	MSpanInuse := float64(memStats.MSpanInuse)
	MSpanSys := float64(memStats.MSpanSys)
	Mallocs := float64(memStats.Mallocs)
	NextGC := float64(memStats.NextGC)
	NumForcedGC := float64(memStats.NumForcedGC)
	NumGC := float64(memStats.NumGC)
	OtherSys := float64(memStats.OtherSys)
	PauseTotalNs := float64(memStats.PauseTotalNs)
	StackInuse := float64(memStats.StackInuse)
	StackSys := float64(memStats.StackSys)
	Sys := float64(memStats.Sys)
	TotalAlloc := float64(memStats.TotalAlloc)
	RandomValue := 1.0 + rand.Float64()*9
	m.Gauges["Alloc"].Value = &Alloc
	m.Gauges["BuckHashSys"].Value = &BuckHashSys
	m.Gauges["Frees"].Value = &Frees
	m.Gauges["GCCPUFraction"].Value = &GCCPUFraction
	m.Gauges["GCSys"].Value = &GCSys
	m.Gauges["HeapAlloc"].Value = &HeapAlloc
	m.Gauges["HeapIdle"].Value = &HeapIdle
	m.Gauges["HeapInuse"].Value = &HeapInuse
	m.Gauges["HeapObjects"].Value = &HeapObjects
	m.Gauges["HeapReleased"].Value = &HeapReleased
	m.Gauges["HeapSys"].Value = &HeapSys
	m.Gauges["LastGC"].Value = &LastGC
	m.Gauges["Lookups"].Value = &Lookups
	m.Gauges["MCacheInuse"].Value = &MCacheInuse
	m.Gauges["MCacheSys"].Value = &MCacheSys
	m.Gauges["MSpanInuse"].Value = &MSpanInuse
	m.Gauges["MSpanSys"].Value = &MSpanSys
	m.Gauges["Mallocs"].Value = &Mallocs
	m.Gauges["NextGC"].Value = &NextGC
	m.Gauges["NumForcedGC"].Value = &NumForcedGC
	m.Gauges["NumGC"].Value = &NumGC
	m.Gauges["OtherSys"].Value = &OtherSys
	m.Gauges["PauseTotalNs"].Value = &PauseTotalNs
	m.Gauges["StackInuse"].Value = &StackInuse
	m.Gauges["StackSys"].Value = &StackSys
	m.Gauges["Sys"].Value = &Sys
	m.Gauges["TotalAlloc"].Value = &TotalAlloc
	m.Gauges["RandomValue"].Value = &RandomValue
	*m.Counters["PollCount"].Delta += 1
	m.mu.Unlock()
}

func (m *MetricsStorage) InitMetrics() {
	m.once.Do(
		func() {
			var memStats runtime.MemStats

			runtime.ReadMemStats(&memStats)
			m.Gauges["Alloc"] = &models.Metrics{ID: "Alloc", MType: "gauge", Value: new(float64)}
			m.Gauges["BuckHashSys"] = &models.Metrics{ID: "BuckHashSys", MType: "gauge", Value: new(float64)}
			m.Gauges["Frees"] = &models.Metrics{ID: "Frees", MType: "gauge", Value: new(float64)}
			m.Gauges["GCCPUFraction"] = &models.Metrics{ID: "GCCPUFraction", MType: "gauge", Value: new(float64)}
			m.Gauges["GCSys"] = &models.Metrics{ID: "GCSys", MType: "gauge", Value: new(float64)}
			m.Gauges["HeapAlloc"] = &models.Metrics{ID: "HeapAlloc", MType: "gauge", Value: new(float64)}
			m.Gauges["HeapIdle"] = &models.Metrics{ID: "HeapIdle", MType: "gauge", Value: new(float64)}
			m.Gauges["HeapInuse"] = &models.Metrics{ID: "HeapInuse", MType: "gauge", Value: new(float64)}
			m.Gauges["HeapObjects"] = &models.Metrics{ID: "HeapObjects", MType: "gauge", Value: new(float64)}
			m.Gauges["HeapReleased"] = &models.Metrics{ID: "HeapReleased", MType: "gauge", Value: new(float64)}
			m.Gauges["HeapSys"] = &models.Metrics{ID: "HeapSys", MType: "gauge", Value: new(float64)}
			m.Gauges["LastGC"] = &models.Metrics{ID: "LastGC", MType: "gauge", Value: new(float64)}
			m.Gauges["Lookups"] = &models.Metrics{ID: "Lookups", MType: "gauge", Value: new(float64)}
			m.Gauges["MCacheInuse"] = &models.Metrics{ID: "MCacheInuse", MType: "gauge", Value: new(float64)}
			m.Gauges["MCacheSys"] = &models.Metrics{ID: "MCacheSys", MType: "gauge", Value: new(float64)}
			m.Gauges["MSpanInuse"] = &models.Metrics{ID: "MSpanInuse", MType: "gauge", Value: new(float64)}
			m.Gauges["MSpanSys"] = &models.Metrics{ID: "MSpanSys", MType: "gauge", Value: new(float64)}
			m.Gauges["Mallocs"] = &models.Metrics{ID: "Mallocs", MType: "gauge", Value: new(float64)}
			m.Gauges["NextGC"] = &models.Metrics{ID: "NextGC", MType: "gauge", Value: new(float64)}
			m.Gauges["NumForcedGC"] = &models.Metrics{ID: "NumForcedGC", MType: "gauge", Value: new(float64)}
			m.Gauges["NumGC"] = &models.Metrics{ID: "NumGC", MType: "gauge", Value: new(float64)}
			m.Gauges["OtherSys"] = &models.Metrics{ID: "OtherSys", MType: "gauge", Value: new(float64)}
			m.Gauges["PauseTotalNs"] = &models.Metrics{ID: "PauseTotalNs", MType: "gauge", Value: new(float64)}
			m.Gauges["StackInuse"] = &models.Metrics{ID: "StackInuse", MType: "gauge", Value: new(float64)}
			m.Gauges["StackSys"] = &models.Metrics{ID: "StackSys", MType: "gauge", Value: new(float64)}
			m.Gauges["Sys"] = &models.Metrics{ID: "Sys", MType: "gauge", Value: new(float64)}
			m.Gauges["TotalAlloc"] = &models.Metrics{ID: "TotalAlloc", MType: "gauge", Value: new(float64)}
			m.Gauges["RandomValue"] = &models.Metrics{ID: "RandomValue", MType: "gauge", Value: new(float64)}
			m.Counters["PollCount"] = &models.Metrics{ID: "PollCount", MType: "counter", Delta: new(int64)}
		})
}
