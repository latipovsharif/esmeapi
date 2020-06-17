package api

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	channel     string = "SMS_CHANNEL"
	logFilePath string = "logs/esmeapi.log"
)

// Server is main server instance
type Server struct {
	l *log.Logger
	e *gin.Engine
	r *rabbit
}

// Run is the entry point to the program
func (s *Server) Run() error {
	lumberjackLogRotate := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    2,   // Max megabytes before log is rotated
		MaxBackups: 500, // Max number of old log files to keep
		MaxAge:     60,  // Max number of days to retain log files
		Compress:   true,
	}
	log.SetOutput(lumberjackLogRotate)

	r, err := newRabbit()
	if err != nil {
		return errors.Wrapf(err, "cannot get rabbit")
	}

	s.r = r
	s.l = log.StandardLogger()
	s.e = gin.Default()
	s.getRoutes()

	if err := s.e.Run(os.Getenv("ESMEAPI_HOST")); err != nil {
		return errors.Wrapf(err, "cannot start server due to %v")
	}

	return nil
}

// Publish message to rabbitmq
func (s *Server) Publish(data []byte) error {
	err := s.r.ch.Publish(
		"",         // exchange
		s.r.q.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})

	if err != nil {
		return errors.Wrap(err, "cannot publish message")
	}

	return nil
}
