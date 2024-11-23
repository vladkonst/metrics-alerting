package storage

import (
	"context"
	"errors"

	"github.com/vladkonst/metrics-alerting/internal/models"
)

type MemStorage struct {
	gauges    map[string]*models.Metrics
	counters  map[string]*models.Metrics
	metricsCh *chan models.Metrics
}

func NewMemStorage(metricsCh *chan models.Metrics) *MemStorage {
	storage := MemStorage{gauges: make(map[string]*models.Metrics), counters: make(map[string]*models.Metrics), metricsCh: metricsCh}
	return &storage
}

func (m *MemStorage) GetCountersValues(ctx context.Context) (map[string]int64, error) {
	countersValues := make(map[string]int64, len(m.counters))

	for k, v := range m.counters {
		countersValues[k] = *v.Delta
	}

	return countersValues, nil
}

func (m *MemStorage) GetGaugesValues(ctx context.Context) (map[string]float64, error) {
	gaugesValues := make(map[string]float64, len(m.gauges))

	for k, v := range m.gauges {
		gaugesValues[k] = *v.Value
	}

	return gaugesValues, nil
}

func (m *MemStorage) AddMetrics(ctx context.Context, metrics []models.Metrics) ([]models.Metrics, error) {
	for i, metric := range metrics {
		switch metric.MType {
		case "counter":
			if _, ok := m.counters[metric.ID]; !ok {
				m.counters[metric.ID] = &metric
			} else {
				*(m.counters[metric.ID].Delta) += *metric.Delta
				metrics[i] = *m.counters[metric.ID]
			}
			metrics[i] = *m.counters[metric.ID]
		case "gauge":
			m.gauges[metric.ID].Value = metric.Value
			metrics[i] = *m.gauges[metric.ID]
		default:
			return nil, errors.New("provided metric type is incorrect")
		}
	}

	return metrics, nil
}

func (m *MemStorage) AddMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, error) {
	switch metric.MType {
	case "counter":
		if _, ok := m.counters[metric.ID]; !ok {
			m.counters[metric.ID] = metric
		} else {
			*(m.counters[metric.ID].Delta) += *metric.Delta
		}
		return m.counters[metric.ID], nil
	case "gauge":
		m.gauges[metric.ID] = metric
		return m.gauges[metric.ID], nil
	default:
		return nil, errors.New("provided metric type is incorrect")
	}
}

func (m *MemStorage) GetMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, error) {
	switch metric.MType {
	case "counter":
		if counter, ok := m.counters[metric.ID]; !ok {
			return nil, errors.New("can't find metric by provided name")
		} else {
			return counter, nil
		}
	case "gauge":
		if gauge, ok := m.gauges[metric.ID]; !ok {
			return nil, errors.New("can't find metric by provided name")
		} else {
			return gauge, nil
		}
	default:
		return nil, errors.New("provided metric type is incorrect")
	}
}
