package agent

import (
	"math/rand/v2"
	"runtime"
	"sync"

	"github.com/vladkonst/metrics-alerting/internal/models"

	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

type MetricsStorage struct {
	once          sync.Once
	RuntimeGauges map[string]*models.Metrics
	PSUtilGauges  map[string]*models.Metrics
	Counters      map[string]*models.Metrics
}

func NewMetricsStorage() (ms MetricsStorage) {
	ms = MetricsStorage{RuntimeGauges: make(map[string]*models.Metrics), PSUtilGauges: make(map[string]*models.Metrics), Counters: make(map[string]*models.Metrics)}
	return
}

func (m *MetricsStorage) UpdateRuntimeMetrics(metricsCh *chan models.Metrics) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	*m.RuntimeGauges["Alloc"].Value = float64(memStats.Alloc)
	*m.RuntimeGauges["BuckHashSys"].Value = float64(memStats.BuckHashSys)
	*m.RuntimeGauges["Frees"].Value = float64(memStats.Frees)
	*m.RuntimeGauges["GCCPUFraction"].Value = float64(memStats.GCCPUFraction)
	*m.RuntimeGauges["GCSys"].Value = float64(memStats.GCSys)
	*m.RuntimeGauges["HeapAlloc"].Value = float64(memStats.HeapAlloc)
	*m.RuntimeGauges["HeapIdle"].Value = float64(memStats.HeapIdle)
	*m.RuntimeGauges["HeapInuse"].Value = float64(memStats.HeapInuse)
	*m.RuntimeGauges["HeapObjects"].Value = float64(memStats.HeapObjects)
	*m.RuntimeGauges["HeapReleased"].Value = float64(memStats.HeapReleased)
	*m.RuntimeGauges["HeapSys"].Value = float64(memStats.HeapSys)
	*m.RuntimeGauges["LastGC"].Value = float64(memStats.LastGC)
	*m.RuntimeGauges["Lookups"].Value = float64(memStats.Lookups)
	*m.RuntimeGauges["MCacheInuse"].Value = float64(memStats.MCacheInuse)
	*m.RuntimeGauges["MCacheSys"].Value = float64(memStats.MCacheSys)
	*m.RuntimeGauges["MSpanInuse"].Value = float64(memStats.MSpanInuse)
	*m.RuntimeGauges["MSpanSys"].Value = float64(memStats.MSpanSys)
	*m.RuntimeGauges["Mallocs"].Value = float64(memStats.Mallocs)
	*m.RuntimeGauges["NextGC"].Value = float64(memStats.NextGC)
	*m.RuntimeGauges["NumForcedGC"].Value = float64(memStats.NumForcedGC)
	*m.RuntimeGauges["NumGC"].Value = float64(memStats.NumGC)
	*m.RuntimeGauges["OtherSys"].Value = float64(memStats.OtherSys)
	*m.RuntimeGauges["PauseTotalNs"].Value = float64(memStats.PauseTotalNs)
	*m.RuntimeGauges["StackInuse"].Value = float64(memStats.StackInuse)
	*m.RuntimeGauges["StackSys"].Value = float64(memStats.StackSys)
	*m.RuntimeGauges["Sys"].Value = float64(memStats.Sys)
	*m.RuntimeGauges["TotalAlloc"].Value = float64(memStats.TotalAlloc)
	*m.RuntimeGauges["RandomValue"].Value = 1.0 + rand.Float64()*9
	*m.Counters["PollCount"].Delta += 1
	for _, v := range m.Counters {
		*metricsCh <- *v
	}

	for _, v := range m.RuntimeGauges {
		*metricsCh <- *v
	}
}

func (m *MetricsStorage) UpdatePSUtilMetrics(metricsCh *chan models.Metrics) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	l, err := load.Avg()
	if err != nil {
		return err
	}

	*m.PSUtilGauges["TotalMemory"].Value = float64(v.Total)
	*m.PSUtilGauges["FreeMemory"].Value = float64(v.Free)
	*m.PSUtilGauges["CPUutilization1"].Value = float64(l.Load1)
	for _, v := range m.PSUtilGauges {
		*metricsCh <- *v
	}

	return nil
}

func (m *MetricsStorage) InitMetrics() {
	m.once.Do(
		func() {
			var memStats runtime.MemStats

			runtime.ReadMemStats(&memStats)
			m.RuntimeGauges["Alloc"] = &models.Metrics{ID: "Alloc", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["BuckHashSys"] = &models.Metrics{ID: "BuckHashSys", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["Frees"] = &models.Metrics{ID: "Frees", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["GCCPUFraction"] = &models.Metrics{ID: "GCCPUFraction", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["GCSys"] = &models.Metrics{ID: "GCSys", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["HeapAlloc"] = &models.Metrics{ID: "HeapAlloc", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["HeapIdle"] = &models.Metrics{ID: "HeapIdle", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["HeapInuse"] = &models.Metrics{ID: "HeapInuse", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["HeapObjects"] = &models.Metrics{ID: "HeapObjects", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["HeapReleased"] = &models.Metrics{ID: "HeapReleased", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["HeapSys"] = &models.Metrics{ID: "HeapSys", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["LastGC"] = &models.Metrics{ID: "LastGC", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["Lookups"] = &models.Metrics{ID: "Lookups", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["MCacheInuse"] = &models.Metrics{ID: "MCacheInuse", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["MCacheSys"] = &models.Metrics{ID: "MCacheSys", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["MSpanInuse"] = &models.Metrics{ID: "MSpanInuse", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["MSpanSys"] = &models.Metrics{ID: "MSpanSys", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["Mallocs"] = &models.Metrics{ID: "Mallocs", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["NextGC"] = &models.Metrics{ID: "NextGC", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["NumForcedGC"] = &models.Metrics{ID: "NumForcedGC", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["NumGC"] = &models.Metrics{ID: "NumGC", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["OtherSys"] = &models.Metrics{ID: "OtherSys", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["PauseTotalNs"] = &models.Metrics{ID: "PauseTotalNs", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["StackInuse"] = &models.Metrics{ID: "StackInuse", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["StackSys"] = &models.Metrics{ID: "StackSys", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["Sys"] = &models.Metrics{ID: "Sys", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["TotalAlloc"] = &models.Metrics{ID: "TotalAlloc", MType: "gauge", Value: new(float64)}
			m.RuntimeGauges["RandomValue"] = &models.Metrics{ID: "RandomValue", MType: "gauge", Value: new(float64)}
			m.PSUtilGauges["TotalMemory"] = &models.Metrics{ID: "TotalMemory", MType: "gauge", Value: new(float64)}
			m.PSUtilGauges["FreeMemory"] = &models.Metrics{ID: "FreeMemory", MType: "gauge", Value: new(float64)}
			m.PSUtilGauges["CPUutilization1"] = &models.Metrics{ID: "CPUutilization1", MType: "gauge", Value: new(float64)}
			m.Counters["PollCount"] = &models.Metrics{ID: "PollCount", MType: "counter", Delta: new(int64)}
		})
}
