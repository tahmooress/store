package logger

import (
	"fmt"
	"io"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
	io.Closer
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.entry.Printf(format, args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.entry.Warningf(format, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logger) Print(args ...interface{}) {
	l.entry.Print(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.entry.Warning(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

func (l *logger) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}

func (l *logger) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}

func (l *logger) Println(args ...interface{}) {
	l.entry.Println(args...)
}

func (l *logger) Warnln(args ...interface{}) {
	l.entry.Warnln(args...)
}

func (l *logger) Warningln(args ...interface{}) {
	l.entry.Warningln(args...)
}

func (l *logger) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}

func (l *logger) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

func (l *logger) Panicln(args ...interface{}) {
	l.entry.Panicln(args...)
}

func (l *logger) Close() error {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return fmt.Errorf("logger: close() >> %w", err)
		}
	}

	return nil
}
