package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/vladkonst/metrics-alerting/internal/models"
)

type FileManager struct {
	filePath      string
	storeInterval int
	metricsCh     *chan models.Metrics
	Metrics       map[string]models.Metrics `json:"metrics"`
}

func NewFileManager(f string, r bool, s int, c *chan models.Metrics) (*FileManager, error) {
	metrics := make(map[string]models.Metrics)
	fm := FileManager{f, s, c, metrics}
	if r {
		if err := fm.InitMetrics(); err != nil {
			return nil, err
		}
	}
	return &fm, nil
}

func (fm *FileManager) InitMetrics() error {
	file, err := os.Open(fm.filePath)
	if err != nil {
		return err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return scanner.Err()
	}

	data := scanner.Bytes()
	if err = json.Unmarshal(data, &fm.Metrics); err != nil {
		return err
	}

	ms := GetStorage(fm.metricsCh)
	for _, metric := range fm.Metrics {
		_, err = ms.AddMetric(&metric)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fm *FileManager) LoadMetricsSync() error {
	for metric := range *fm.metricsCh {
		fm.Metrics[metric.ID] = metric
		file, err := os.Create(fm.filePath)
		if err != nil {
			return err
		}

		defer file.Close()
		enc := json.NewEncoder(file)
		if err := enc.Encode(fm.Metrics); err != nil {
			return err
		}
	}
	return nil
}

func (fm *FileManager) LoadMetrics() error {
	if fm.storeInterval == 0 {
		go fm.LoadMetricsSync()
		return nil
	}

	tc := time.NewTicker(time.Second * time.Duration(fm.storeInterval))
	for {
		select {
		case <-tc.C:
			file, err := os.Create(fm.filePath)
			if err != nil {
				return err
			}

			defer file.Close()
			enc := json.NewEncoder(file)
			if err := enc.Encode(fm.Metrics); err != nil {
				return err
			}

		case metric := <-*fm.metricsCh:
			fm.Metrics[metric.ID] = metric
		}
	}
}
