package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vladkonst/metrics-alerting/handlers"
)

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", handlers.GetMetricsPage)

	r.Route("/value", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			{
				http.Error(w, "Bad request.", http.StatusBadRequest)
			}
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Route("/{type}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Metric not found.", http.StatusNotFound)
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			})
			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			})
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			})
			r.Get("/{name}", handlers.GetCurrentMetricValue)
		})
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Bad request.", http.StatusBadRequest) })
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Route("/{type}", func(r chi.Router) {
			r.Post("/", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Metric not found.", http.StatusNotFound) })
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			})
			r.Put("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			})
			r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			})
			r.Route("/{name}", func(r chi.Router) {
				r.Post("/", func(w http.ResponseWriter, r *http.Request) { http.Error(w, "Metric not found.", http.StatusNotFound) })
				r.Post("/{value}", handlers.UpdateMetric)
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
				})
				r.Put("/", func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
				})
				r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
				})
			})
		})
	})

	return r
}
