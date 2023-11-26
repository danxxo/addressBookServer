package errorlogger

import (
	"log"
	"os"
)

type ErrorLogger struct {
	file *os.File
}

func NewErrorLogger(filename string) (*ErrorLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &ErrorLogger{file}, nil
}

func (el *ErrorLogger) LogError(err error) {
	log.SetOutput(el.file)
	log.Println(err)
}

func (el *ErrorLogger) Close() {
	el.file.Close()
}
