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
	fmt.Println(" Initializing Zerolog :: Service Name : ", os.Getenv(("SERVICE_NAME")))
	file, err := os.OpenFile(
		fmt.Sprintf(`%s.log`, os.Getenv("SERVICE_NAME")),
		os.O_CREATE|os.O_APPEND|os.O_RDWR|os.O_TRUNC,
		0666,
	)

	if err != nil {
		fmt.Sprintln(" Initializing Zerolog :: Error while opening log file, Error : ", err)
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
