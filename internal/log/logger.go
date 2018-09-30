package log

import (
	"context"
	"github.com/sirupsen/logrus"
)

type loggerKeyType int

const loggerKey loggerKeyType = iota

var loggerEntry *logrus.Entry

func init() {
	logger := logrus.New()
	loggerEntry = logger.WithField("reqID", "main")
}

func NewContext(ctx context.Context, fields logrus.Fields) context.Context {
	ctxLogger := WithCtx(ctx).WithFields(fields)
	return context.WithValue(ctx, loggerKey, ctxLogger)
}

func WithCtx(ctx context.Context) *logrus.Entry {

	if ctx == nil {
		return loggerEntry
	}

	if ctxLogger, ok := ctx.Value(loggerKey).(*logrus.Entry); ok {
		return ctxLogger
	}

	return loggerEntry
}
