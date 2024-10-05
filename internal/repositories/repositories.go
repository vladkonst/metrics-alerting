package repositories

type GaugeRepository interface {
	AddGauge(name string, value float64) error
	GetGauge(name string) (float64, error)
	RemoveGauge(name string) error
}

type CounterRepository interface {
	AddCounter(name string, value int64) error
	GetCounter(name string) (int64, error)
	RemoveCounter(name string) error
}
