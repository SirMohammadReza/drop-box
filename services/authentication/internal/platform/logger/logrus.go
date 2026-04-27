package logger

import "github.com/sirupsen/logrus"

type LogrusWrapper struct {
	entry *logrus.Entry
}

func NewLogrusWrapper(l *logrus.Logger) *LogrusWrapper {
	return &LogrusWrapper{
		entry: logrus.NewEntry(l),
	}
}

func (lw *LogrusWrapper) Errorf(format string, args ...any) {
	lw.entry.Errorf(format, args...)
}

func (lw *LogrusWrapper) Infof(format string, args ...any) {
	lw.entry.Infof(format, args...)
}

func (lw *LogrusWrapper) WithField(key string, value any) Logger {
	return &LogrusWrapper{
		entry: lw.entry.WithField(key, value),
	}
}
