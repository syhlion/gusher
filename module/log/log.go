package log

import (
	"github.com/Sirupsen/logrus"
)

var (
	Logger *logrus.Logger = nil
)

func init() {

	if Logger == nil {
		Logger = logrus.New()
	}
}
