package app

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vladkonst/metrics-alerting/handlers"
	"github.com/vladkonst/metrics-alerting/internal/configs"
	"github.com/vladkonst/metrics-alerting/internal/models"
	"github.com/vladkonst/metrics-alerting/internal/storage"
)

type App struct {
	Storage         handlers.MetricRepository
	MetricsChan     *chan models.Metrics
	StorageProvider *handlers.StorageProvider
	done            *chan bool
	cfg             *configs.ServerCfg
}

func NewApp(done *chan bool) *App {
	cfg := configs.GetServerConfig()
	ps := cfg.IntervalsCfg.DatabaseDSN
	var s handlers.MetricRepository
	var conn *sql.DB
	metricsCh := make(chan models.Metrics)
	s = storage.NewMemStorage(&metricsCh)
	if ps != "" {
		conn, _ = sql.Open("pgx", ps)
	}
	// switch ps {
	// case "":
	// 	s = storage.NewMemStorage(&metricsCh)
	// default:
	// 	conn, err := sql.Open("pgx", ps)
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	s = storage.NewPGStorage(conn)
	// }

	sp := &handlers.StorageProvider{Storage: s, MetricsChan: &metricsCh, DB: conn}
	return &App{Storage: s, MetricsChan: &metricsCh, StorageProvider: sp, done: done, cfg: cfg}
}

func (a *App) GetMetricsChanel() *chan models.Metrics {
	return a.MetricsChan
}

func (a *App) GetStorage() handlers.MetricRepository {
	return a.Storage
}

func (a App) Run() {
	fileStorage, err := storage.NewFileManager(a.cfg.IntervalsCfg.FileStoragePath, a.cfg.IntervalsCfg.Restore, a.cfg.IntervalsCfg.StoreInterval, a.MetricsChan, a.Storage)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		if err := fileStorage.ProcessMetrics(); err != nil {
			log.Panic(err)
		}
	}()

	go func() {
		log.Panic(http.ListenAndServe(a.cfg.NetAddressCfg.String(), a.GetRouter()))
	}()

	<-*a.done
	fileStorage.LoadMetrics()
}

func (a *App) GetRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", a.StorageProvider.GetMetricsPage)

	r.Get("/ping", a.StorageProvider.PingDB)

	r.Route("/value", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			{
				http.Error(w, "Bad request.", http.StatusBadRequest)
			}
		})
		r.Post("/", a.StorageProvider.GetMetric)
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
			r.Get("/{name}", a.StorageProvider.GetGaugeMetricValue)
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
			r.Get("/{name}", a.StorageProvider.GetCounterMetricValue)
		})
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
		})
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/", a.StorageProvider.UpdateMetric)
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
				r.Post("/{value}", a.StorageProvider.UpdateGaugeMetric)
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
				r.Post("/{value}", a.StorageProvider.UpdateCounterMetric)
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
