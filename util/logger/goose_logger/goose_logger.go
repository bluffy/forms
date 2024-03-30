package goose_logger

import (
	log "github.com/sirupsen/logrus"
)

type logger struct {
}

func New() *logger {
	return &logger{}
}

func (l *logger) Fatal(args ...interface{}) {
	log.Fatal(args...)
}
func (l *logger) Fatalf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}
func (l *logger) Print(args ...interface{}) {
	log.Info(args...)
}
func (l *logger) Printf(format string, v ...interface{}) {
	log.Infof(format, v...)
}
func (l *logger) Println(args ...interface{}) {
	log.Infoln(args...)
}
