package logger

func Error(s string) {
	Logger.Error().Msgf(s)
}

func Debug(s string) {
	Logger.Debug().Msgf(s)
}

func Info(s string) {
	Logger.Info().Msgf(s)
}

func Trace(s string) {
	Logger.Trace().Msgf(s)
}

func Warn(s string) {
	Logger.Warn().Msgf(s)
}

func Panic(s string) {
	Logger.Panic().Msgf(s)
}

func Fatal(s string) {
	Logger.Fatal().Msgf(s)
}
