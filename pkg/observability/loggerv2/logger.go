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

	file, err := os.OpenFile(
		fmt.Sprintf(`%s.log`, serviceName),
		os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC,
		0666,
	)
	if err != nil {
		fmt.Sprintln("Initializing Zerolog :: Error while opening log file, Error : ", err)
		return err
	}

	zerolog.TimeFieldFormat = time.RFC3339

	var writer io.Writer = file

	if os.Getenv("ENABLE_LOGGERV2") == "true" {
		writer = zerolog.MultiLevelWriter(file, zerolog.ConsoleWriter{
			Out: os.Stderr,
			FieldsExclude: []string{
				"user_agent",
				// "span_id",
				// "trace_id",
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

	defaultWriter := zerolog.ConsoleWriter{
		Out: os.Stderr,
	}

	defaultLogger := zerolog.New(defaultWriter).
		With().
		Timestamp().
		Caller().
		Logger()

	zerolog.DefaultContextLogger = &defaultLogger

	return nil
}
