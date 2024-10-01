package storage

type Gauge struct {
	name  string
	value float64
}

type Counter struct {
	name  string
	value int64
}

type MemStorage struct {
	gauges   map[string]Gauge
	counters map[string]Counter
}

var storage MemStorage

func GetStorage() *MemStorage {
	if storage.gauges == nil {
		storage = MemStorage{gauges: make(map[string]Gauge), counters: make(map[string]Counter)}
	}
	return &storage
}

func (m *MemStorage) AddCounter(name string, value int64) error {
	counter := m.counters[name]
	counter.value = value
	m.counters[name] = counter
	return nil
}

func (m *MemStorage) GetCounter(name string) (Counter, error) {
	return m.counters[name], nil
}

func (m *MemStorage) RemoveCounter(name string) error {
	delete(m.counters, name)
	return nil
}

func (m *MemStorage) AddGauge(name string, value float64) error {
	gauge := m.gauges[name]
	gauge.value = value
	m.gauges[name] = gauge
	return nil
}

func (m *MemStorage) GetGauge(name string) (Gauge, error) {
	return m.gauges[name], nil
}

func (m *MemStorage) RemoveGauge(name string) error {
	delete(m.gauges, name)
	return nil
}
