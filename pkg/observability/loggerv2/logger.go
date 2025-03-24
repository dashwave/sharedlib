package loggerv2

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"

)

var ll zerolog.Logger
var once sync.Once

func Get() zerolog.Logger {
	once.Do(func() {
		if err := initLogger(); err != nil {
			panic(err)
		}
	})
	return ll
}

func Ctx(ctx context.Context) *zerolog.Logger {
	return log.Ctx(ctx)
}

func ZCtx(ctx context.Context) zerolog.Context {
	return Ctx(ctx).With()
}

func initLogger() error {
	
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		return fmt.Errorf("SERVICE_NAME is not set")
	}

	fmt.Println("Initializing Zerolog :: Service Name :: ", serviceName)

	logFile := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s.log", serviceName),
		MaxSize:    1 * 1000, // 1 GB
		MaxBackups: 3,
		MaxAge:     7,    // days
		Compress:   true, // gzip old logs
	}

	var writer io.Writer = logFile

	if os.Getenv("ENABLE_LOGGERV2") == "true" {
		writer = zerolog.MultiLevelWriter(logFile, zerolog.ConsoleWriter{
			Out: os.Stderr,
			FieldsExclude: []string{
				"user_agent",
			},
		})
	}

	ll = zerolog.
		New(writer).
		With().
		Str("service", serviceName).
		Timestamp().
		Caller().
		Logger()

	defaultWriter := zerolog.ConsoleWriter{Out: os.Stderr}
	defaultLogger := zerolog.New(defaultWriter).With().Timestamp().Caller().Logger()
	zerolog.DefaultContextLogger = &defaultLogger


	return nil
}
