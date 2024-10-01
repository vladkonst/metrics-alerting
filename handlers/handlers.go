package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/vladkonst/metrics-alerting/internal/repositories"
	"github.com/vladkonst/metrics-alerting/internal/storage"
)

func UpdateGauge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	pathValues := strings.Split(strings.Trim(r.URL.Path, "/"), "/gauge/")

	if len(pathValues) == 1 {
		http.Error(w, "Metric not found.", http.StatusNotFound)
		return
	}

	metricData := strings.Split(pathValues[1], "/")

	if len(metricData) != 2 {
		http.Error(w, "Metric not found.", http.StatusNotFound)
		return
	}

	v, err := strconv.ParseFloat(metricData[1], 64)

	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}
	var memStorage repositories.GaugeRepository = storage.GetStorage()
	memStorage.AddGauge(metricData[0], v)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
}

func UpdateCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathValues := strings.Split(strings.Trim(r.URL.Path, "/"), "/counter/")

	if len(pathValues) == 1 {
		http.Error(w, "Metric not found.", http.StatusNotFound)
		return
	}

	metricData := strings.Split(pathValues[1], "/")

	if len(metricData) != 2 {
		http.Error(w, "Metric not found.", http.StatusNotFound)
		return
	}

	v, err := strconv.ParseInt(metricData[1], 10, 64)

	if err != nil {
		http.Error(w, "Bad request.", http.StatusBadRequest)
		return
	}

	var memStorage repositories.CounterRepository = storage.GetStorage()
	memStorage.AddCounter(metricData[0], v)

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
}
