package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/0987363/2table-backend/models"

	"net/url"
	"time"
)

var (
	logLevel string
	logStore string //日志存储地方
	host     string
)

func LoggerInit() *logrus.Logger {
	logger := logrus.New()
	//host, _ := os.Hostname()

	logger.Level = models.ConvertLevel(logLevel)
	logger.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: time.RFC3339Nano}

	return logger
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := url.QueryUnescape(c.Request.URL.String())
		if err != nil {
			u = c.Request.URL.String()
		}

		log := LoggerInit().WithFields(logrus.Fields{
			RequestIDKey: GetRequestID(c),
			"Method":     c.Request.Method,
			"UserAgent":  c.Request.UserAgent(),
			"Url":        u,
			"Host":       c.Request.Host,
		})

		start := time.Now()
		c.Set(models.MiddwareKeyLogger, log)
		c.Next()

		if c.Request.Method == "OPTIONS" || u == "/version" {
			return
		}
		spent := int(time.Now().Sub(start) / time.Millisecond)
		if spent > 1000 {
			log.WithFields(logrus.Fields{
				"Status": c.Writer.Status(),
				"Spent":  spent,
			}).Warnf("Responded %03d in %dms", c.Writer.Status(), spent)
		} else {
			log.WithFields(logrus.Fields{
				"Status": c.Writer.Status(),
				"Spent":  spent,
			}).Infof("Responded %03d in %dms", c.Writer.Status(), spent)
		}
	}
}

func GetLogger(c *gin.Context) *logrus.Entry {
	if logger, ok := c.Get(models.MiddwareKeyLogger); ok {
		return logger.(*logrus.Entry)
	}

	return nil
}

func SetLoggerField(c *gin.Context, k, v string) {
	if logger, ok := c.Get(models.MiddwareKeyLogger); ok {
		c.Set(models.MiddwareKeyLogger, logger.(*logrus.Entry).WithField(k, v))
	}
}
