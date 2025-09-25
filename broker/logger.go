package broker

import (
	"log"
	"os"
	"path/filepath"
)

type Logger interface {
	Enabled() bool
	Info(level int, msg string, keysAndValues ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
	WithValues(keysAndValues ...interface{}) Logger
	WithName(name string) Logger
}
type DefaultLogger struct {
	logger *log.Logger
	name   string
	values []interface{}
	level  int
}

func NewDefaultLogger(filePath string) (*DefaultLogger, error) {
	dir := filepath.Dir(filePath)
	if ex := os.MkdirAll(dir, 0755); ex != nil {
		return nil, ex
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &DefaultLogger{
		logger: log.New(f, "", log.LstdFlags|log.Lshortfile),
	}, nil
}

func (l *DefaultLogger) Enabled() bool {
	// Có thể mở rộng thành check level
	return true
}

func (l *DefaultLogger) Info(level int, msg string, keysAndValues ...interface{}) {
	if !l.Enabled() {
		return
	}
	l.logger.Printf("[INFO] %s %s %v", l.name, msg, append(l.values, keysAndValues...))
}

func (l *DefaultLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if !l.Enabled() {
		return
	}
	l.logger.Printf("[ERROR] %s %s err=%v %v", l.name, msg, err, append(l.values, keysAndValues...))
}

func (l *DefaultLogger) WithValues(keysAndValues ...interface{}) Logger {
	// copy để giữ immutability
	newLogger := *l
	newLogger.values = append(newLogger.values, keysAndValues...)
	return &newLogger
}

func (l *DefaultLogger) WithName(name string) Logger {
	newLogger := *l
	if l.name != "" {
		newLogger.name = l.name + "." + name
	} else {
		newLogger.name = name
	}
	return &newLogger
}
