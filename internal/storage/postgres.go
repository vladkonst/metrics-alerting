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
	return &PGStorage{conn: conn}
}

func (s PGStorage) Bootstrap(ctx context.Context) error {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// tx.ExecContext(ctx, `
	//     CREATE TABLE users (
	//         id varchar(128) PRIMARY KEY,
	//         username varchar(128)
	//     )
	// `)
	// tx.ExecContext(ctx, `CREATE UNIQUE INDEX sender_idx ON users (username)`)

	// tx.ExecContext(ctx, `
	//     CREATE TABLE messages (
	//         id serial PRIMARY KEY,
	//         sender varchar(128),
	//         recipient varchar(128),
	//         payload text,
	//         sent_at timestamp with time zone,
	//         read_at timestamp with time zone DEFAULT NULL
	//     )
	// `)
	// tx.ExecContext(ctx, `CREATE INDEX recipient_idx ON messages (recipient)`)

	return tx.Commit()
}

func (s *PGStorage) GetCountersValues() (map[string]int64, error) {

	return nil, nil
}

func (s *PGStorage) GetGaugesValues() (map[string]float64, error) {
	return nil, nil
}

func (s *PGStorage) AddMetric(metric *models.Metrics) (*models.Metrics, error) {
	switch metric.MType {
	case "counter":

	case "gauge":

	default:
		return nil, errors.New("provided metric type is incorrect")
	}
	return nil, nil
}

func (s *PGStorage) GetMetric(metric *models.Metrics) (*models.Metrics, error) {
	switch metric.MType {
	case "counter":

	case "gauge":

	default:
		return nil, errors.New("provided metric type is incorrect")
	}
	return nil, nil
}
