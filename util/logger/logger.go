package logger

import (
	"context"
	"fmt"
	"os"
	"worker-manager/middleware"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func logPrefix(ctx context.Context) string {
	if ctx == nil {
		return "[No Context] "
	}
	reqId, ok := middleware.GetRequestId(ctx)
	if !ok {
		return "[No Request ID] "
	}
	return fmt.Sprintf("[%s] ", reqId)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logrus.Infof(logPrefix(ctx)+format, args...)
}

func Infoln(ctx context.Context, args ...interface{}) {
	logrus.Infoln(append([]interface{}{logPrefix(ctx)}, args...)...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logrus.Debugf(logPrefix(ctx)+format, args...)
}

func Debugln(ctx context.Context, args ...interface{}) {
	logrus.Debugln(append([]interface{}{logPrefix(ctx)}, args...)...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logrus.Errorf(logPrefix(ctx)+format, args...)
}

func Errorln(ctx context.Context, args ...interface{}) {
	logrus.Errorln(append([]interface{}{logPrefix(ctx)}, args...)...)
}
