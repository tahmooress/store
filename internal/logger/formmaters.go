package logger

import (
	"time"

	"github.com/sirupsen/logrus"
)

// nolint: gochecknoglobals
var defaultTextFormmater = &logrus.TextFormatter{
	FullTimestamp:   true,
	TimestampFormat: time.RFC3339,
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime:  "@Timestamp",
		logrus.FieldKeyMsg:   "@Message",
		logrus.FieldKeyLevel: "@Level",
		logrus.FieldKeyFile:  "@File",
		logrus.FieldKeyFunc:  "@Func",
	},
}

// nolint: gochecknoglobals
var defaultJSONFormmater = &logrus.JSONFormatter{
	DisableTimestamp: false,
	TimestampFormat:  time.RFC3339,
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime:  "timestamp",
		logrus.FieldKeyMsg:   "message",
		logrus.FieldKeyLevel: "level",
		logrus.FieldKeyFile:  "file",
		logrus.FieldKeyFunc:  "func",
	},
	PrettyPrint: true,
}
