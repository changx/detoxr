package ds

import (
	"os"

	"github.com/gobuffalo/logger"
	"github.com/sirupsen/logrus"
)

func JSONLogger(lvl logger.Level) logger.FieldLogger {
	l := logrus.New()
	l.Level = lvl
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetOutput(os.Stdout)
	return logger.Logrus{FieldLogger: l}
}
