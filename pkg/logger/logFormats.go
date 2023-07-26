package logger

import (
	"fmt"
	"strings"
)

func Error(a ...interface{}) {
	format := strings.Repeat("%v ", len(a))
	fmt.Println(a)
	Logger.Error().Msgf(format, a...)
}

func Debug(a ...interface{}) {
	format := strings.Repeat("%v ", len(a))
	fmt.Println(a)
	Logger.Debug().Msgf(format, a...)
}

func Info(a ...interface{}) {
	format := strings.Repeat("%v ", len(a))
	fmt.Println(a)
	Logger.Info().Msgf(format, a...)
}

func Trace(a ...interface{}) {
	format := strings.Repeat("%v ", len(a))
	fmt.Println(a)
	Logger.Trace().Msgf(format, a...)
}

func Warn(a ...interface{}) {
	format := strings.Repeat("%v ", len(a))
	fmt.Println(a)
	Logger.Warn().Msgf(format, a...)
}

func Panic(a ...interface{}) {
	format := strings.Repeat("%v ", len(a))
	fmt.Println(format, a)
	Logger.Panic().Msgf(format, a...)
}

func Fatal(a ...interface{}) {
	format := strings.Repeat("%v ", len(a))
	fmt.Println(a)
	Logger.Fatal().Msgf(format, a...)
}
