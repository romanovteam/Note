package logger

import (
	"log"
	"os"
	"path/filepath"
)

type Logger interface {
	LogError(err error)
}

// FileLogger Реализация интерфейса для записи логов в файл
type FileLogger struct {
	file *os.File
}

// NewFileLogger создает новый FileLogger
func NewFileLogger(logDir string) (*FileLogger, error) {
	logFilePath := filepath.Join(logDir, "app.log")
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

// LogError записывает сообщение об ошибке в лог
func (l *FileLogger) LogError(err error) {
	log.SetOutput(l.file)
	log.Println("ERROR", err)
}

func (l *FileLogger) Close() {
	err := l.file.Close()
	if err != nil {
		return
	}
}
