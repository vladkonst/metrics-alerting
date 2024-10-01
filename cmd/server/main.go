package main

import (
	"net/http"

	"github.com/vladkonst/metrics-alerting/handlers"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/update/gauge/", handlers.UpdateGauge)

	mux.HandleFunc("/update/counter/", handlers.UpdateCounter)

	mux.HandleFunc("/update/", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Bad request.", http.StatusBadRequest) })

	http.ListenAndServe("localhost:8080", mux)
}
