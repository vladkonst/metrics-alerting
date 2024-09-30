package repositories

import "github.com/vladkonst/metrics-alerting/internal/storage"

type GaugeRepository interface {
	AddGauge(name string, value float64) error
	GetGauge(name string) (storage.Gauge, error)
	RemoveGauge(name string) error
}

type CounterRepository interface {
	AddCounter(name string, value int64) error
	GetCounter(name string) (storage.Counter, error)
	RemoveCounter(name string) error
}
