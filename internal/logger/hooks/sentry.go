package hooks

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func NewSentryHook(dsn, env string, sampleRate float64, levels []logrus.Level) (logrus.Hook, error) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         dsn,
		Environment: env,
		SampleRate:  sampleRate,
	})
	if err != nil {
		return nil, fmt.Errorf("hook: NewSentryHook >> %w", err)
	}

	levelMap := map[logrus.Level]sentry.Level{
		logrus.TraceLevel: sentry.LevelDebug,
		logrus.DebugLevel: sentry.LevelDebug,
		logrus.InfoLevel:  sentry.LevelInfo,
		logrus.WarnLevel:  sentry.LevelWarning,
		logrus.ErrorLevel: sentry.LevelError,
		logrus.FatalLevel: sentry.LevelFatal,
		logrus.PanicLevel: sentry.LevelFatal,
	}

	return &sentryHook{
		hub:    sentry.CurrentHub(),
		levels: levels,
		lvlMap: levelMap,
	}, nil
}

type sentryHook struct {
	hub    *sentry.Hub
	levels []logrus.Level
	lvlMap map[logrus.Level]sentry.Level
}

func (s *sentryHook) Fire(entry *logrus.Entry) error {
	s.hub.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(s.lvlMap[entry.Level])
		scope.SetExtras(entry.Data)

		if ok, err := getErrorFromEntry(entry); ok {
			scope.SetExtra("log.message", entry.Message)
			s.hub.CaptureException(err)
		} else {
			s.hub.CaptureMessage(entry.Message)
		}
	})

	return nil
}

func getErrorFromEntry(entry *logrus.Entry) (bool, error) {
	err, ok := entry.Data[logrus.ErrorKey].(error)

	return ok, err
}

func (s *sentryHook) Levels() []logrus.Level {
	return s.levels
}
