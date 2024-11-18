package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vladkonst/metrics-alerting/internal/logger"
	"github.com/vladkonst/metrics-alerting/internal/models"
)

func GzipMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w
		if strings.Contains(r.Header.Get("Content-Type"), "application/json") || strings.Contains(r.Header.Get("Accept"), "text/html") {
			acceptEncoding := r.Header.Get("Accept-Encoding")
			supportsGzip := strings.Contains(acceptEncoding, "gzip")
			if supportsGzip {
				cw := newCompressWriter(w)
				ow = cw
				defer cw.Close()
				cw.Header().Add("Content-Encoding", "gzip")
			}

			contentEncoding := r.Header.Get("Content-Encoding")
			sendsGzip := strings.Contains(contentEncoding, "gzip")
			if sendsGzip {
				cr, err := newCompressReader(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				r.Body = cr
				defer cr.Close()
			}
		}
		h.ServeHTTP(ow, r)
	})
}

type loggingResponseWriter struct {
	r      http.ResponseWriter
	status int
	size   int
}

func (lr *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lr.r.Write(b)
	lr.size += size
	return size, err
}

func (lr *loggingResponseWriter) WriteHeader(statusCode int) {
	lr.r.WriteHeader(statusCode)
	lr.status = statusCode
}

func (lr *loggingResponseWriter) Header() http.Header {
	return lr.r.Header()
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logger := logger.Get()
		lw := loggingResponseWriter{r: w}
		next.ServeHTTP(&lw, r)
		logger.
			Info().
			Str("method", r.Method).
			Str("URI", r.URL.RequestURI()).
			Dur("duration", time.Since(start)).
			Int("status", lw.status).
			Int("size", lw.size).
			Msg("incoming request")
	})
}

type MetricRepository interface {
	AddMetric(*models.Metrics) (*models.Metrics, error)
	GetMetric(*models.Metrics) (*models.Metrics, error)
	GetGaugesValues() (map[string]float64, error)
	GetCountersValues() (map[string]int64, error)
	GetMetricsChanel() *chan models.Metrics
}

type StorageProvider struct {
	handler func(http.ResponseWriter, *http.Request, MetricRepository)
	storage MetricRepository
}

func NewStorageProvider(handlerToWrap func(http.ResponseWriter, *http.Request, MetricRepository), memStorage MetricRepository) http.HandlerFunc {
	sp := &StorageProvider{handlerToWrap, memStorage}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { sp.handler(w, r, sp.storage) })
}

func GetMetric(w http.ResponseWriter, r *http.Request, memStorage MetricRepository) {
	metric := new(models.Metrics)
	dec := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	metric, err := memStorage.GetMetric(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateMetric(w http.ResponseWriter, r *http.Request, memStorage MetricRepository) {
	metricsCh := memStorage.GetMetricsChanel()
	metric := new(models.Metrics)
	dec := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metric, err := memStorage.AddMetric(metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	*metricsCh <- *metric
}

func GetMetricsPage(w http.ResponseWriter, r *http.Request, memStorage MetricRepository) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	gauges, err := memStorage.GetGaugesValues()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	counters, err := memStorage.GetCountersValues()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Gauges   map[string]float64
		Counters map[string]int64
	}{
		Gauges:   gauges,
		Counters: counters,
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
}

func GetGaugeMetricValue(w http.ResponseWriter, r *http.Request, memStorage MetricRepository) {
	metric := models.Metrics{ID: chi.URLParam(r, "name"), MType: "gauge"}
	gauge, err := memStorage.GetMetric(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, fmt.Sprintf("%g", *gauge.Value))
}

func GetCounterMetricValue(w http.ResponseWriter, r *http.Request, memStorage MetricRepository) {
	metric := models.Metrics{ID: chi.URLParam(r, "name"), MType: "counter"}
	counter, err := memStorage.GetMetric(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, fmt.Sprintf("%d", *counter.Delta))
}

func UpdateGaugeMetric(w http.ResponseWriter, r *http.Request, memStorage MetricRepository) {
	metricsCh := memStorage.GetMetricsChanel()
	v, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	metric := models.Metrics{ID: chi.URLParam(r, "name"), Value: &v, MType: "gauge"}
	_, err = memStorage.AddMetric(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	*metricsCh <- metric
}

func UpdateCounterMetric(w http.ResponseWriter, r *http.Request, memStorage MetricRepository) {
	metricsCh := memStorage.GetMetricsChanel()
	v, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	metric := models.Metrics{ID: chi.URLParam(r, "name"), Delta: &v, MType: "counter"}
	_, err = memStorage.AddMetric(&metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	*metricsCh <- metric
}
