package loggerv2

import (
	"github.com/rs/zerolog"
)

func Trace() *zerolog.Event {
	return l.Trace()
}

func Debug() *zerolog.Event {
	return l.Debug()
}

func Info() *zerolog.Event {
	return l.Info()
}

func Warn() *zerolog.Event {
	return l.Warn()
}

func Error() *zerolog.Event {
	return l.Error()
}

func Fatal() *zerolog.Event {
	return l.Fatal()
}

func Panic() *zerolog.Event {
	return l.Panic()
}
