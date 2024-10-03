package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vladkonst/metrics-alerting/internal/repositories"
	"github.com/vladkonst/metrics-alerting/internal/storage"
)

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

func GetCurrentMetricValue(w http.ResponseWriter, r *http.Request) {
	switch chi.URLParam(r, "type") {
	case "gauge":
		{
			var memStorage repositories.GaugeRepository = storage.GetStorage()
			gauge, err := memStorage.GetGauge(chi.URLParam(r, "name"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			io.WriteString(w, fmt.Sprintf("%f", gauge))
			w.WriteHeader(http.StatusOK)
		}
	case "counter":
		{
			var memStorage repositories.CounterRepository = storage.GetStorage()
			counter, err := memStorage.GetCounter(chi.URLParam(r, "name"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			io.WriteString(w, fmt.Sprintf("%d", counter))
			w.WriteHeader(http.StatusOK)
		}
	default:
		{
			http.Error(w, "Bad request.", http.StatusBadRequest)
			return
		}
	}
}

func UpdateMetric(w http.ResponseWriter, r *http.Request) {
	switch chi.URLParam(r, "type") {
	case "gauge":
		{
			v, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
			if err != nil {
				http.Error(w, "Bad request.", http.StatusBadRequest)
				return
			}
			var memStorage repositories.GaugeRepository = storage.GetStorage()
			memStorage.AddGauge(chi.URLParam(r, "name"), v)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
		}
	case "counter":
		{
			v, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
			if err != nil {
				http.Error(w, "Bad request.", http.StatusBadRequest)
				return
			}
			var memStorage repositories.CounterRepository = storage.GetStorage()
			memStorage.AddCounter(chi.URLParam(r, "name"), v)
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
		}
	default:
		{
			http.Error(w, "Bad request.", http.StatusBadRequest)
			return
		}
	}
}
