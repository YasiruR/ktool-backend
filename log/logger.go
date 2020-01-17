package log

import (
	"github.com/pickme-go/log"
)

var Logger log.Logger

func Init() {
	logLevel := log.Level(Cfg.Level)
	Logger = log.Constructor.Log(log.WithColors(Cfg.Colors), log.WithLevel(logLevel), log.WithFilePath(Cfg.FilePathEnabled), log.Prefixed(`level-1`))
}