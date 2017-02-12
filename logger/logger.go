package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/hooks/test"
	"github.com/spf13/viper"
)

var log *logrus.Logger
var hook *test.Hook

func GetLogger() *logrus.Logger {
	if viper.GetBool("plis.testing") {
		if log == nil && hook == nil {
			log, hook = test.NewNullLogger()
		}
		return log
	}
	if log == nil {
		log = logrus.New()
		return log
	}
	return log
}

func SetLevel(level logrus.Level) {
	if log == nil {
		log = logrus.New()
	}
	log.Level = level
}
func GetHook() *test.Hook {
	if log == nil && hook == nil {
		log, hook = test.NewNullLogger()
	}
	return hook
}
