package loggerv2

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func InitLogger() (*zerolog.Logger, error) {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		return nil, fmt.Errorf("SERVICE_NAME is not set")
	}
	fmt.Println("Initializing Zerolog :: Service Name :: ", serviceName)

	file, err := os.OpenFile(
		fmt.Sprintf(`%s.log`, serviceName),
		os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC,
		0666,
	)
	if err != nil {
		fmt.Sprintln("Initializing Zerolog :: Error while opening log file, Error : ", err)
		return nil, err
	}
	defer file.Close()

	zerolog.TimeFieldFormat = time.RFC3339

	var writer io.Writer = file

	if os.Getenv("ENV") == "development" {
		writer = zerolog.MultiLevelWriter(file, zerolog.ConsoleWriter{Out: os.Stderr})
	}

	l := zerolog.
		New(writer).
		With().
		Str("service", serviceName).
		Timestamp().
		Logger()

	return &l, nil
}
