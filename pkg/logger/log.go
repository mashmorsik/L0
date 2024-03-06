package logger

import (
	"github.com/rs/zerolog"
	"os"
	"time"
)

var Nlog *zerolog.Logger

func BuildLogger() {
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()

	Nlog = &logger
}

func Errf(format string, v ...interface{}) {
	Nlog.Error().Msgf(format, v...)
}

func Err(err error, msg string) {
	Nlog.Error().Err(err).Msg(msg)
}

func Warn(msg string) {
	Nlog.Warn().Msg(msg)
}

func Infof(format string, v ...interface{}) {
	Nlog.Info().Msgf(format, v...)
}

func Info(msg string) {
	Nlog.Info().Msg(msg)
}

func Fatalf(format string, v ...interface{}) {
	Nlog.Fatal().Msgf(format, v...)
}
