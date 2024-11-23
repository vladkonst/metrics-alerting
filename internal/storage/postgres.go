package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/vladkonst/metrics-alerting/internal/models"
)

type PGStorage struct {
	conn *sql.DB
}

func NewPGStorage(conn *sql.DB) *PGStorage {
	storage := PGStorage{conn: conn}
	storage.Bootstrap(context.Background())
	return &storage
}

func (s PGStorage) Bootstrap(ctx context.Context) error {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	tx.ExecContext(ctx, `
	    CREATE TABLE counters (
	        name varchar PRIMARY KEY,
			value bigint
	    )
	`)
	tx.ExecContext(ctx, `
	    CREATE TABLE gauges (
	        name varchar PRIMARY KEY,
	        value double precision
	    )
	`)
	return tx.Commit()
}

func (s *PGStorage) GetCountersValues(ctx context.Context) (map[string]int64, error) {
	counters := make(map[string]int64)
	rows, err := s.conn.QueryContext(ctx, "SELECT name, value FROM counters")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		var value int64
		err = rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}

		counters[name] = value
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return counters, nil
}

func (s *PGStorage) GetGaugesValues(ctx context.Context) (map[string]float64, error) {
	gauges := make(map[string]float64)
	rows, err := s.conn.QueryContext(ctx, "SELECT name, value FROM gauges")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		var value float64
		err = rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}

		gauges[name] = value
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return gauges, nil
}

func (s *PGStorage) AddMetrics(ctx context.Context, metrics []models.Metrics) ([]models.Metrics, error) {
	addedCounters := make(map[string]int64)
	addedGauges := make(map[string]bool)
	tx, err := s.conn.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	for i, metric := range metrics {
		switch metric.MType {
		case "counter":
			var counterValue int64
			row := s.conn.QueryRowContext(ctx, `SELECT value  FROM counters WHERE name = $1`, metric.ID)
			err := row.Scan(&counterValue)
			if _, ok := addedCounters[metric.ID]; err != nil && !ok {
				if _, err := tx.ExecContext(ctx, "INSERT INTO counters (name, value) VALUES($1,$2)", metric.ID, *metric.Delta); err != nil {
					return nil, err
				}
				addedCounters[metric.ID] += *metric.Delta
			} else {
				*metric.Delta = counterValue + *metric.Delta
				if _, err := tx.ExecContext(ctx, "UPDATE counters SET value=$1 WHERE name=$2", addedCounters[metric.ID]+*metric.Delta, metric.ID); err != nil {
					return nil, err
				}
			}
			metrics[i] = metric
		case "gauge":
			var gaugeName string
			row := s.conn.QueryRowContext(ctx, `SELECT name  FROM counters WHERE name = $1`, metric.ID)
			err := row.Scan(&gaugeName)
			if _, ok := addedGauges[metric.ID]; err != nil && !ok {
				if _, err := tx.ExecContext(ctx, "INSERT INTO gauges (name, value) VALUES($1,$2)", metric.ID, *metric.Value); err != nil {
					return nil, err
				}
				addedGauges[metric.ID] = true
			} else {
				if _, err := tx.ExecContext(ctx, "UPDATE gauges SET value=$1 WHERE name=$2", *metric.Value, metric.ID); err != nil {
					return nil, err
				}
			}

		default:
			return nil, errors.New("provided metric type is incorrect")
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return metrics, nil
}

func (s *PGStorage) AddMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, error) {
	switch metric.MType {
	case "counter":
		var counterValue int64
		row := s.conn.QueryRowContext(ctx, `SELECT value  FROM counters WHERE name = $1`, metric.ID)
		err := row.Scan(&counterValue)
		if err != nil {
			if _, err := s.conn.ExecContext(ctx, "INSERT INTO counters (name, value) VALUES($1,$2)", metric.ID, metric.Delta); err != nil {
				return nil, err
			}
		} else {
			if _, err := s.conn.ExecContext(ctx, "UPDATE counters SET value=$1 where name=$2", *metric.Delta+counterValue, metric.ID); err != nil {
				return nil, err
			}
		}
		return metric, nil
	case "gauge":
		var gaugeValue float64
		row := s.conn.QueryRowContext(ctx, `SELECT value  FROM gauges WHERE name = $1`, metric.ID)
		err := row.Scan(&gaugeValue)
		if err != nil {
			if _, err := s.conn.ExecContext(ctx, "INSERT INTO gauges (name, value) VALUES($1,$2)", metric.ID, metric.Value); err != nil {
				return nil, err
			}
		} else {
			metric.Value = &gaugeValue
			if _, err := s.conn.ExecContext(ctx, "UPDATE gauges SET value=$1 where name=$2", *metric.Value, metric.ID); err != nil {
				return nil, err
			}
		}
		return metric, nil
	default:
		return nil, errors.New("provided metric type is incorrect")
	}
}

func (s *PGStorage) GetMetric(ctx context.Context, metric *models.Metrics) (*models.Metrics, error) {
	switch metric.MType {
	case "counter":
		row := s.conn.QueryRowContext(ctx, `SELECT value  FROM counters WHERE name = $1`, metric.ID)
		err := row.Scan(&metric.Delta)
		if err != nil {
			return nil, errors.New("can't find metric by provided name")
		}
		return metric, nil
	case "gauge":
		row := s.conn.QueryRowContext(ctx, `SELECT value  FROM gauges WHERE name = $1`, metric.ID)
		err := row.Scan(&metric.Value)
		if err != nil {
			return nil, errors.New("can't find metric by provided name")
		}
		return metric, nil
	default:
		return nil, errors.New("provided metric type is incorrect")
	}
}
