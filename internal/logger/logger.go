package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var (
	once sync.Once
	log  zerolog.Logger
)

func Get() zerolog.Logger {
	once.Do(func() {
		zerolog.TimeFieldFormat = time.RFC3339Nano
		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		log = zerolog.New(output).
			Level(zerolog.InfoLevel).
			With().
			Timestamp().
			Logger()
	})

	return log
}
