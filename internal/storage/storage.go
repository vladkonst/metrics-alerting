package storage

import (
	"errors"

	"github.com/vladkonst/metrics-alerting/internal/models"
)

type MemStorage struct {
	Gauges   map[string]*models.Metrics
	Counters map[string]*models.Metrics
}

var storage MemStorage

func (m *MemStorage) GetCountersValues() (map[string]int64, error) {
	countersValues := make(map[string]int64, len(m.Counters))

	for k, v := range m.Counters {
		countersValues[k] = v.Delta
	}

	return countersValues, nil
}

func (m *MemStorage) GetGaugesValues() (map[string]float64, error) {
	gaugesValues := make(map[string]float64, len(m.Gauges))

	for k, v := range m.Gauges {
		gaugesValues[k] = v.Value
	}

	return gaugesValues, nil
}

func GetStorage() *MemStorage {
	if storage.Gauges == nil {
		storage = MemStorage{Gauges: make(map[string]*models.Metrics), Counters: make(map[string]*models.Metrics)}
	}
	return &storage
}

func (m *MemStorage) AddMetric(metric *models.Metrics) (*models.Metrics, error) {
	switch metric.MType {
	case "counter":
		if _, ok := m.Counters[metric.ID]; !ok {
			m.Counters[metric.ID] = metric
		} else {
			m.Counters[metric.ID].Delta += metric.Delta
		}
		return m.Counters[metric.ID], nil
	case "gauge":
		m.Gauges[metric.ID] = metric
		return m.Gauges[metric.ID], nil
	default:
		return nil, errors.New("provided metric type is incorrect")
	}
}

func (m *MemStorage) GetMetric(metric *models.Metrics) (*models.Metrics, error) {
	switch metric.MType {
	case "counter":
		counter, ok := m.Counters[metric.ID]
		if !ok {
			return metric, nil
		}
		return counter, nil
	case "gauge":
		gauge, ok := m.Gauges[metric.ID]
		if !ok {
			return metric, nil
		}
		return gauge, nil
	default:
		return nil, errors.New("provided metric type is incorrect")
	}
}
