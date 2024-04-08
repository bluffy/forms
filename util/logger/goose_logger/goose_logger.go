package goose_logger

import "github.com/sirupsen/logrus"

type logger struct {
}

func New() *logger {
	return &logger{}
}

func (l *logger) Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}
func (l *logger) Fatalf(format string, v ...interface{}) {
	logrus.Fatalf(format, v...)
}
func (l *logger) Print(args ...interface{}) {
	logrus.Info(args...)
}
func (l *logger) Printf(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}
func (l *logger) Println(args ...interface{}) {
	logrus.Infoln(args...)
}
