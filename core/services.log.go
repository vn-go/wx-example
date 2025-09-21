package core

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type logService interface {
}
type logLumberjack struct {
}

func newLogLumberjack(cfg *configInfo) (logService, error) {
	//make folder if not exist the folfer path is in cfg.Path
	if err := os.MkdirAll(cfg.Log.Path, os.ModeAppend); err != nil {
		return nil, err
	}
	logFileName := filepath.Join(cfg.Log.Path, "app.log")
	log.SetOutput(&lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    cfg.Log.Size,     // MB, roll up if exeed
		MaxBackups: cfg.Log.Backup,   // nu, of file need to be backup
		MaxAge:     cfg.Log.Age,      // older than {cfg.Log.Age} days delete
		Compress:   cfg.Log.Compress, // compress old files
	})
	return &logLumberjack{}, nil
}
