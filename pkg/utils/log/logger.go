package log

import (
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

type Logger struct {
	Level      log.Level
	SlogLogger slog.Logger
	Logger     log.Logger
}

func CreateLogFile(fileName string) (*os.File, error) {
	logFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer logFile.Close()

	return logFile, nil
}

func NewLogger(level log.Level, prefix, prefixColor string) *Logger {
	// logFile, err := CreateLogFile("log.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// multiWriter := io.MultiWriter(os.Stdout, logFile)

	styles := log.DefaultStyles()
	styles.Prefix = lipgloss.NewStyle().
		Foreground(lipgloss.Color(prefixColor)).
		Bold(true)

	handler := log.NewWithOptions(os.Stdout, log.Options{
		Prefix:          prefix,
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		Level:           level,
		Formatter:       log.TextFormatter,
	})
	handler.SetStyles(styles)

	slogLogger := slog.New(handler)

	return &Logger{
		Level:      level,
		SlogLogger: *slogLogger,
		Logger:     *handler,
	}
}

func (logger *Logger) Log(message string) {
	logger.Logger.Print(message)
}

func (logger *Logger) LogInfo(message string) {
	logger.Logger.Info(message)
}

func (logger *Logger) LogSimulation(logValue slog.Value, message string) {
	logger.Logger.Print(message, "Sim", logValue)
}

func (logger *Logger) Error(message string, err error) {
	logger.Logger.Error(message, "Error", err)
}

func (logger *Logger) Fatal(message string, err error) {
	logger.Logger.Fatal(message, "Error", err)
}

func (logger *Logger) Debug(message string) {
	logger.Logger.Debug(message)
}
