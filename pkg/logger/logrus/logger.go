package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func New(cfg Config) (*logrus.Logger, io.Closer, error) {
	log := logrus.New()
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
		log.Warnf("Некорректный уровень логирования 'level'. Используется значение по умолчанию - 'info'")
	}
	log.SetLevel(level)
	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     false,
		})
	}
	var closer io.Closer
	if cfg.OutputFile != "" {
		file, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, nil, fmt.Errorf("не удалось открыть файл лога: %w", err)
		}
		log.SetOutput(file)
		closer = file
	} else {
		log.SetOutput(os.Stdout)
		closer = nil
	}

	return log, closer, nil
}
