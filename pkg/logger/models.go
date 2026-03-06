package logger

import (
	"io"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Fatal() *zerolog.Event
	Panic() *zerolog.Event
	With() zerolog.Context
	Level(level zerolog.Level) Logger
	Output(w io.Writer) Logger
}

type loggerWrapper struct {
	zerolog.Logger
}
