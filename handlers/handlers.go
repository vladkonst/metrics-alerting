package handlers

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
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

type Hasher struct {
	hash hash.Hash
}

func NewHasher(key string) *Hasher {
	if key == "" {
		return nil
	}
	hash := sha256.New()
	return &Hasher{hash}
}

func (h *Hasher) HashBody(b []byte) (string, error) {
	_, err := h.hash.Write(b)
	if err != nil {
		return "", errors.New("internal server error")
	}

	dst := h.hash.Sum(nil)
	dstHex := hex.EncodeToString(dst)
	return dstHex, nil
}

func (h *Hasher) HashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		dst, _ := h.HashBody(b)
		src := r.Header.Get("HashSHA256")
		if src != dst {
			http.Error(w, "invalid hash provided", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(b))
		next.ServeHTTP(w, r)
	})
}

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
	AddMetrics(context.Context, []models.Metrics) ([]models.Metrics, error)
	AddMetric(context.Context, *models.Metrics) (*models.Metrics, error)
	GetMetric(context.Context, *models.Metrics) (*models.Metrics, error)
	GetGaugesValues(context.Context) (map[string]float64, error)
	GetCountersValues(context.Context) (map[string]int64, error)
}

type StorageProvider struct {
	Storage     MetricRepository
	DB          *sql.DB
	MetricsChan *chan models.Metrics
}

func (sp *StorageProvider) PingDB(w http.ResponseWriter, r *http.Request) {
	if sp.DB == nil {
		http.Error(w, "database connection failed", http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := sp.DB.PingContext(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (sp *StorageProvider) GetMetric(w http.ResponseWriter, r *http.Request) {
	metric := new(models.Metrics)
	dec := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	metric, err := sp.Storage.GetMetric(ctx, metric)
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

func (sp *StorageProvider) UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := make([]models.Metrics, 0)
	dec := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(&metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	metrics, err := sp.Storage.AddMetrics(ctx, metrics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(metrics); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, metric := range metrics {
		*sp.MetricsChan <- metric
	}
}

func (sp *StorageProvider) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric := new(models.Metrics)
	dec := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")
	if err := dec.Decode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	metric, err := sp.Storage.AddMetric(ctx, metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(metric); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	*sp.MetricsChan <- *metric
}

func (sp *StorageProvider) GetMetricsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx, cancel := context.WithTimeout(r.Context(), 6*time.Second)
	defer cancel()
	gauges, err := sp.Storage.GetGaugesValues(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	counters, err := sp.Storage.GetCountersValues(ctx)
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

func (sp *StorageProvider) GetGaugeMetricValue(w http.ResponseWriter, r *http.Request) {
	metric := models.Metrics{ID: chi.URLParam(r, "name"), MType: "gauge"}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	gauge, err := sp.Storage.GetMetric(ctx, &metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, fmt.Sprintf("%g", *gauge.Value))
}

func (sp *StorageProvider) GetCounterMetricValue(w http.ResponseWriter, r *http.Request) {
	metric := models.Metrics{ID: chi.URLParam(r, "name"), MType: "counter"}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	counter, err := sp.Storage.GetMetric(ctx, &metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, fmt.Sprintf("%d", *counter.Delta))
}

func (sp *StorageProvider) UpdateGaugeMetric(w http.ResponseWriter, r *http.Request) {
	metricsCh := sp.MetricsChan
	v, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	metric := models.Metrics{ID: chi.URLParam(r, "name"), Value: &v, MType: "gauge"}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	_, err = sp.Storage.AddMetric(ctx, &metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	*metricsCh <- metric
}

func (sp *StorageProvider) UpdateCounterMetric(w http.ResponseWriter, r *http.Request) {
	metricsCh := sp.MetricsChan
	v, err := strconv.ParseInt(chi.URLParam(r, "value"), 10, 64)
	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	metric := models.Metrics{ID: chi.URLParam(r, "name"), Delta: &v, MType: "counter"}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	_, err = sp.Storage.AddMetric(ctx, &metric)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	*metricsCh <- metric
}
