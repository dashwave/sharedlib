package logger

func Error(format string, a ...interface{}) {
	Logger.Error().Msgf(format, a...)
}

func Debug(format string, a ...interface{}) {
	Logger.Debug().Msgf(format, a...)
}

func Info(format string, a ...interface{}) {
	Logger.Info().Msgf(format, a...)
}

func Trace(format string, a ...interface{}) {
	Logger.Trace().Msgf(format, a...)
}

func Warn(format string, a ...interface{}) {
	Logger.Warn().Msgf(format, a...)
}

func Panic(format string, a ...interface{}) {
	Logger.Panic().Msgf(format, a...)
}

func Fatal(format string, a ...interface{}) {
	Logger.Fatal().Msgf(format, a...)
}
