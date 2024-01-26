package logger

import (
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var once sync.Once

var log zerolog.Logger

func Get() zerolog.Logger {
	once.Do(func() {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

		zerolog.SetGlobalLevel(zerolog.InfoLevel)

		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

		log = zerolog.New(output).With().Timestamp().Logger()

	})

	return log
}
