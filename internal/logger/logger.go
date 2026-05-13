package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func New(env, level string) *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	if strings.EqualFold(env, "production") {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	}

	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		parsedLevel = logrus.InfoLevel
	}
	log.SetLevel(parsedLevel)

	return log
}
