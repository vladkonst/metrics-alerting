package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/vladkonst/metrics-alerting/handlers"
	"github.com/vladkonst/metrics-alerting/internal/storage"
)

func GetRouter() http.Handler {
	r := chi.NewRouter()
	memStorage := storage.GetStorage(nil)

	r.Get("/", handlers.NewStorageProvider(handlers.GetMetricsPage, memStorage))

	r.Route("/value", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			{
				http.Error(w, "Bad request.", http.StatusBadRequest)
			}
		})
		r.Post("/", handlers.NewStorageProvider(handlers.GetMetric, memStorage))
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Route("/gauge", func(r chi.Router) {
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
			r.Get("/{name}", handlers.NewStorageProvider(handlers.GetGaugeMetricValue, memStorage))
		})
		r.Route("/counter", func(r chi.Router) {
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
			r.Get("/{name}", handlers.NewStorageProvider(handlers.GetCounterMetricValue, memStorage))
		})
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
		})
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", handlers.NewStorageProvider(handlers.UpdateMetric, memStorage))
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		})
		r.Route("/gauge", func(r chi.Router) {
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
				r.Post("/{value}", handlers.NewStorageProvider(handlers.UpdateGaugeMetric, memStorage))
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
		r.Route("/counter", func(r chi.Router) {
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
				r.Post("/{value}", handlers.NewStorageProvider(handlers.UpdateCounterMetric, memStorage))
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
		r.Post("/*", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
		})
	})

	return handlers.GzipMiddleware(handlers.LogRequest(r))
}
