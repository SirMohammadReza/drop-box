package logger

type Logger interface {
	Errorf(format string, args ...any)
	Infof(format string, args ...any)
	WithField(key string, value any) Logger
}
