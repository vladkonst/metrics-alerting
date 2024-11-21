package storage

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/vladkonst/metrics-alerting/handlers"
	"github.com/vladkonst/metrics-alerting/internal/models"
)

type FileManager struct {
	filePath      string
	storeInterval int
	metricsCh     *chan models.Metrics
	Metrics       map[string]models.Metrics `json:"metrics"`
}

func NewFileManager(f string, r bool, s int, c *chan models.Metrics, storage handlers.MetricRepository) (*FileManager, error) {
	metrics := make(map[string]models.Metrics)
	fm := FileManager{f, s, c, metrics}
	if r {
		if err := fm.InitMetrics(storage); err != nil {
			return nil, err
		}
	}
	return &fm, nil
}

func (fm *FileManager) InitMetrics(storage handlers.MetricRepository) error {
	file, err := os.OpenFile(fm.filePath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	dec := json.NewDecoder(file)
	if err = dec.Decode(&fm.Metrics); err != nil && err.Error() != "EOF" {
		return err
	}

	for _, metric := range fm.Metrics {
		_, err = storage.AddMetric(context.Background(), &metric)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileManager) LoadMetrics() error {
	file, err := os.OpenFile(fm.filePath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	defer file.Close()
	enc := json.NewEncoder(file)
	if err := enc.Encode(fm.Metrics); err != nil {
		return err
	}
	return nil
}

func (fm *FileManager) ProcessMetricsSync() error {
	for metric := range *fm.metricsCh {
		fm.Metrics[metric.ID] = metric
		if err := fm.LoadMetrics(); err != nil {
			return err
		}
	}

	return nil
}

func (fm *FileManager) ProcessMetrics() error {
	if fm.storeInterval == 0 {
		go fm.ProcessMetricsSync()
		return nil
	}

	tc := time.NewTicker(time.Second * time.Duration(fm.storeInterval))
	for {
		select {
		case <-tc.C:
			if err := fm.LoadMetrics(); err != nil {
				return err
			}
		case metric := <-*fm.metricsCh:
			fm.Metrics[metric.ID] = metric
		}
	}
}
