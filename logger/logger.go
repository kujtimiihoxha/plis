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
	} else {
		if log == nil {
			log = logrus.New()
			return log
		} else {
			return log
		}
	}
}
func GetHook() *test.Hook {
	if log == nil && hook == nil {
		log, hook = test.NewNullLogger()
	}
	return hook
}
