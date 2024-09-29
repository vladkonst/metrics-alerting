package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

var memStorage = MemStorage{map[string]float64{}, map[string]int64{}}

func updateGauge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	pathValues := strings.Split(strings.Trim(r.URL.Path, "/"), "/gauge/")

	if len(pathValues) == 1 {
		http.Error(w, "Bad request.", http.StatusBadRequest)
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

	memStorage.gauges[metricData[0]] = v
}

func updateCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathValues := strings.Split(strings.Trim(r.URL.Path, "/"), "/counter/")

	if len(pathValues) == 1 {
		http.Error(w, "Bad request.", http.StatusBadRequest)
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

	memStorage.counters[metricData[0]] += v
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/gauge/", updateGauge)

	mux.HandleFunc("/update/counter/", updateCounter)

	mux.HandleFunc("/update/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Bad request.", http.StatusBadRequest) }))

	http.ListenAndServe("localhost:8080", mux)
}
