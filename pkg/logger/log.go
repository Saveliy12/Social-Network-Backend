package logger

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var log LoggerInterface

// LoggerInterface определяет интерфейс для логирования
type LoggerInterface interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

// Logger реализует LoggerInterface
type Logger struct {
	*logrus.Logger
}

// CustomFormatter реализует logrus.Formatter интерфейс
type CustomFormatter struct{}

// Format реализует метод Format для кастомного форматирования логов
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	_, filename, line, _ := runtime.Caller(8)
	scriptName := filepath.Base(filename)

	message := fmt.Sprintf("[%s] [%s] [%s:%d] %s\n",
		entry.Time.Format("2006-01-02 15:04:05"), // Дата-время
		entry.Level.String(),                     // Уровень логирования
		scriptName,                               // Имя файла
		line,                                     // Номер строки
		entry.Message,                            // Сообщение
	)

	return []byte(message), nil
}

func (l *Logger) Debug(msg string) {
	l.Debug(msg)
}

func (l *Logger) Info(msg string) {
	l.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.Error(msg)
}

func (l *Logger) Fatal(msg string) {
	l.Fatal(msg)
}

// InitLogger инициализирует глобальный логгер
func InitLogger() {
	logger := logrus.New()

	// Включение отчета о вызове для включения имени файла и номера строки в лог
	logger.SetReportCaller(true)

	// Установка пользовательского форматтера
	logger.SetFormatter(&CustomFormatter{})

	// Установка уровня логирования по умолчанию
	logger.SetLevel(logrus.DebugLevel)

	log = &Logger{logger}
}

// GetLogger возвращает глобальный логгер
func GetLogger() LoggerInterface {
	return log
}
