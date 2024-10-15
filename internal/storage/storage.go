package storage

import (
	"errors"
)

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]int64
}

var storage MemStorage

func GetStorage() *MemStorage {
	if storage.Gauges == nil {
		storage = MemStorage{Gauges: make(map[string]float64), Counters: make(map[string]int64)}
	}
	return &storage
}

func (m *MemStorage) AddCounter(name string, value int64) error {
	m.Counters[name] += value
	return nil
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
	if counter, ok := m.Counters[name]; !ok {
		return 0, errors.New("can't find metric by provided name")
	} else {
		return counter, nil
	}
}

func (m *MemStorage) RemoveCounter(name string) error {
	delete(m.Counters, name)
	return nil
}

func (m *MemStorage) AddGauge(name string, value float64) error {
	m.Gauges[name] = value
	return nil
}

func (m *MemStorage) GetGauge(name string) (float64, error) {
	if gauge, ok := m.Gauges[name]; !ok {
		return 0.0, errors.New("can't find metric by provided name")
	} else {
		return gauge, nil
	}
}

func (m *MemStorage) RemoveGauge(name string) error {
	delete(m.Gauges, name)
	return nil
}
