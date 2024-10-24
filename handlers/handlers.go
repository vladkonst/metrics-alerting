package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/vladkonst/metrics-alerting/internal/storage"
)

type GaugeRepository interface {
	AddGauge(name string, value float64) error
	GetGauge(name string) (float64, error)
}

type CounterRepository interface {
	AddCounter(name string, value int64) error
	GetCounter(name string) (int64, error)
}

type GaugeStorageProvider struct {
	handler    func(http.ResponseWriter, *http.Request, GaugeRepository)
	memStorage GaugeRepository
}

func NewGaugeStorageProvider(handlerToWrap func(http.ResponseWriter, *http.Request, GaugeRepository)) *GaugeStorageProvider {
	var memStorage GaugeRepository = storage.GetStorage()
	return &GaugeStorageProvider{handlerToWrap, memStorage}
}

func (sp *GaugeStorageProvider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sp.handler(w, r, sp.memStorage)
}

type CounterStorageProvider struct {
	handler    func(http.ResponseWriter, *http.Request, CounterRepository)
	memStorage CounterRepository
}

func NewCounterStorageProvider(handlerToWrap func(http.ResponseWriter, *http.Request, CounterRepository)) *CounterStorageProvider {
	var memStorage CounterRepository = storage.GetStorage()
	return &CounterStorageProvider{handlerToWrap, memStorage}
}

func (sp *CounterStorageProvider) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sp.handler(w, r, sp.memStorage)
}

func GetMetricsPage(w http.ResponseWriter, r *http.Request) {
	memStorage := storage.GetStorage()

	data := struct {
		Gauges   map[string]float64
		Counters map[string]int64
	}{
		Gauges:   memStorage.Gauges,
		Counters: memStorage.Counters,
	}

	tmpl := `
	<!DOCTYPE html>
	<head>
		<meta charset="UTF-8">
		<title>Metrics list</title>
	</head>
	<body>
		<ul>
		{{range $key, $value := .Gauges}}
			<li>{{$key}}: {{$value}}</li>
		{{end}}
		{{range $key, $value := .Counters}}
			<li>{{$key}}: {{$value}}</li>
		{{end}}
		</ul>
	</body>
	</html>`

	t, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "text/html; charset=utf-8")
}

func GetGaugeMetricValue(w http.ResponseWriter, r *http.Request, memStorage GaugeRepository) {
	gauge, err := memStorage.GetGauge(chi.URLParam(r, "name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, fmt.Sprintf("%g", gauge))
}

func GetCounterMetricValue(w http.ResponseWriter, r *http.Request, memStorage CounterRepository) {
	counter, err := memStorage.GetCounter(chi.URLParam(r, "name"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, fmt.Sprintf("%d", counter))
}

func UpdateGaugeMetric(w http.ResponseWriter, r *http.Request, memStorage GaugeRepository) {
	v, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}
	memStorage.AddGauge(chi.URLParam(r, "name"), v)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
}

func UpdateCounterMetric(w http.ResponseWriter, r *http.Request, memStorage CounterRepository) {
	v, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}
	memStorage.AddCounter(chi.URLParam(r, "name"), v)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
}
