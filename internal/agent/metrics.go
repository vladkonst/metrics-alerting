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

	m.Gauges["Alloc"].Value = float64(memStats.Alloc)
	m.Gauges["BuckHashSys"].Value = float64(memStats.BuckHashSys)
	m.Gauges["Frees"].Value = float64(memStats.Frees)
	m.Gauges["GCCPUFraction"].Value = float64(memStats.GCCPUFraction)
	m.Gauges["GCSys"].Value = float64(memStats.GCSys)
	m.Gauges["HeapAlloc"].Value = float64(memStats.HeapAlloc)
	m.Gauges["HeapIdle"].Value = float64(memStats.HeapIdle)
	m.Gauges["HeapInuse"].Value = float64(memStats.HeapInuse)
	m.Gauges["HeapObjects"].Value = float64(memStats.HeapObjects)
	m.Gauges["HeapReleased"].Value = float64(memStats.HeapReleased)
	m.Gauges["HeapSys"].Value = float64(memStats.HeapSys)
	m.Gauges["LastGC"].Value = float64(memStats.LastGC)
	m.Gauges["Lookups"].Value = float64(memStats.Lookups)
	m.Gauges["MCacheInuse"].Value = float64(memStats.MCacheInuse)
	m.Gauges["MCacheSys"].Value = float64(memStats.MCacheSys)
	m.Gauges["MSpanInuse"].Value = float64(memStats.MSpanInuse)
	m.Gauges["MSpanSys"].Value = float64(memStats.MSpanSys)
	m.Gauges["Mallocs"].Value = float64(memStats.Mallocs)
	m.Gauges["NextGC"].Value = float64(memStats.NextGC)
	m.Gauges["NumForcedGC"].Value = float64(memStats.NumForcedGC)
	m.Gauges["NumGC"].Value = float64(memStats.NumGC)
	m.Gauges["OtherSys"].Value = float64(memStats.OtherSys)
	m.Gauges["PauseTotalNs"].Value = float64(memStats.PauseTotalNs)
	m.Gauges["StackInuse"].Value = float64(memStats.StackInuse)
	m.Gauges["StackSys"].Value = float64(memStats.StackSys)
	m.Gauges["Sys"].Value = float64(memStats.Sys)
	m.Gauges["TotalAlloc"].Value = float64(memStats.TotalAlloc)
	m.Gauges["RandomValue"].Value = 1.0 + rand.Float64()*9
	m.Counters["PollCount"].Delta++
	m.mu.Unlock()
}

func (m *MetricsStorage) InitMetrics() {
	m.once.Do(
		func() {
			var memStats runtime.MemStats

			runtime.ReadMemStats(&memStats)
			m.Gauges["Alloc"] = &models.Metrics{ID: "Alloc", MType: "gauge", Value: float64(memStats.Alloc)}
			m.Gauges["BuckHashSys"] = &models.Metrics{ID: "BuckHashSys", MType: "gauge", Value: float64(memStats.BuckHashSys)}
			m.Gauges["Frees"] = &models.Metrics{ID: "Frees", MType: "gauge", Value: float64(memStats.Frees)}
			m.Gauges["GCCPUFraction"] = &models.Metrics{ID: "GCCPUFraction", MType: "gauge", Value: float64(memStats.GCCPUFraction)}
			m.Gauges["GCSys"] = &models.Metrics{ID: "GCSys", MType: "gauge", Value: float64(memStats.GCSys)}
			m.Gauges["HeapAlloc"] = &models.Metrics{ID: "HeapAlloc", MType: "gauge", Value: float64(memStats.HeapAlloc)}
			m.Gauges["HeapIdle"] = &models.Metrics{ID: "HeapIdle", MType: "gauge", Value: float64(memStats.HeapIdle)}
			m.Gauges["HeapInuse"] = &models.Metrics{ID: "HeapInuse", MType: "gauge", Value: float64(memStats.HeapInuse)}
			m.Gauges["HeapObjects"] = &models.Metrics{ID: "HeapObjects", MType: "gauge", Value: float64(memStats.HeapObjects)}
			m.Gauges["HeapReleased"] = &models.Metrics{ID: "HeapReleased", MType: "gauge", Value: float64(memStats.HeapReleased)}
			m.Gauges["HeapSys"] = &models.Metrics{ID: "HeapSys", MType: "gauge", Value: float64(memStats.HeapSys)}
			m.Gauges["LastGC"] = &models.Metrics{ID: "LastGC", MType: "gauge", Value: float64(memStats.LastGC)}
			m.Gauges["Lookups"] = &models.Metrics{ID: "Lookups", MType: "gauge", Value: float64(memStats.Lookups)}
			m.Gauges["MCacheInuse"] = &models.Metrics{ID: "MCacheInuse", MType: "gauge", Value: float64(memStats.MCacheInuse)}
			m.Gauges["MCacheSys"] = &models.Metrics{ID: "MCacheSys", MType: "gauge", Value: float64(memStats.MCacheSys)}
			m.Gauges["MSpanInuse"] = &models.Metrics{ID: "MSpanInuse", MType: "gauge", Value: float64(memStats.MSpanInuse)}
			m.Gauges["MSpanSys"] = &models.Metrics{ID: "MSpanSys", MType: "gauge", Value: float64(memStats.MSpanSys)}
			m.Gauges["Mallocs"] = &models.Metrics{ID: "Mallocs", MType: "gauge", Value: float64(memStats.Mallocs)}
			m.Gauges["NextGC"] = &models.Metrics{ID: "NextGC", MType: "gauge", Value: float64(memStats.NextGC)}
			m.Gauges["NumForcedGC"] = &models.Metrics{ID: "NumForcedGC", MType: "gauge", Value: float64(memStats.NumForcedGC)}
			m.Gauges["NumGC"] = &models.Metrics{ID: "NumGC", MType: "gauge", Value: float64(memStats.NumGC)}
			m.Gauges["OtherSys"] = &models.Metrics{ID: "OtherSys", MType: "gauge", Value: float64(memStats.OtherSys)}
			m.Gauges["PauseTotalNs"] = &models.Metrics{ID: "PauseTotalNs", MType: "gauge", Value: float64(memStats.PauseTotalNs)}
			m.Gauges["StackInuse"] = &models.Metrics{ID: "StackInuse", MType: "gauge", Value: float64(memStats.StackInuse)}
			m.Gauges["StackSys"] = &models.Metrics{ID: "StackSys", MType: "gauge", Value: float64(memStats.StackSys)}
			m.Gauges["Sys"] = &models.Metrics{ID: "Sys", MType: "gauge", Value: float64(memStats.Sys)}
			m.Gauges["TotalAlloc"] = &models.Metrics{ID: "TotalAlloc", MType: "gauge", Value: float64(memStats.TotalAlloc)}
			m.Gauges["RandomValue"] = &models.Metrics{ID: "RandomValue", MType: "gauge", Value: 1.0 + rand.Float64()*9}
			m.Counters["PollCount"] = &models.Metrics{ID: "PollCount", MType: "counter", Delta: 0}
		})
}
