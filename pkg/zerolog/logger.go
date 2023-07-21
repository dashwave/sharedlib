package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var (
	Logger zerolog.Logger
)

func init() {
	file, err := os.OpenFile(
		fmt.Sprintf(`%s.log`, os.Getenv("SERVICE_NAME")),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)

	if err != nil {
		panic(err)
	}
	//defer file.Close()

	zerolog.TimeFieldFormat = time.RFC3339

	l := zerolog.
		New(file).
		With().
		Timestamp().
		Logger()
	Logger = l
}
