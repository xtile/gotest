package arbilogger

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type ArbiLogger struct {
	config *Config
	logger *logrus.Logger
}

func New(s *Config) *ArbiLogger {
	return &ArbiLogger{
		config: s,
		logger: logrus.New(),
	}
}

func (s *ArbiLogger) Start() error {

	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Info("starting API server")

	return http.ListenAndServe(s.config.BindAddr, s.router)

}

func (s *ArbiLogger) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)

	if err != nil {
		return err
	}
	s.logger.SetLevel(level)

	return nil
}
