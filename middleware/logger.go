package middleware

import (
	"fmt"
	"io/ioutil"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/gin-gonic/gin"
	logrusloki "github.com/hongyuyang/logrus-loki"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/0987363/2table-backend/models"

	"net"
	"net/url"
	"os"
	"time"
)

const (
	loggerStoreElk  = "elk"  //elk
	loggerStoreLoki = "loki" //loki
)

var (
	logConn    net.Conn         //elk
	lokiClient *logrusloki.Loki //loki client
	logLevel   string
	logStore   string //日志存储地方
	host       string
)

func ConnectLogger(level, dst, store string) error {
	logLevel = level
	logStore = store
	host, _ = os.Hostname()

	if dst != "" {
		if logStore == loggerStoreElk || logStore == "" {
			conn, err := net.Dial("udp", dst)
			if err != nil {
				return err
			}
			logConn = conn
		} else if logStore == loggerStoreLoki {
			//loki客户端会开启一个协程批量发送数据, 不用每次写日志初始化
			var err error
			lokiClient, err = newLokiClient(dst)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func LoggerInit() *logrus.Logger {
	logger := logrus.New()
	//host, _ := os.Hostname()

	logger.Level = models.ConvertLevel(logLevel)
	logger.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true, TimestampFormat: time.RFC3339Nano}
	if logStore == loggerStoreElk || logStore == "" {
		//elk
		if logConn == nil {
			return logger
		}

		fields := GetLoggerCommonFields()
		fields["Hostname"] = host
		fields["Release"] = viper.GetString("release")
		hook := logrustash.New(logConn, logrustash.LogstashFormatter{
			Fields: fields,
			Formatter: &logrus.JSONFormatter{
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime: "@timestamp",
					logrus.FieldKeyMsg:  "message",
				},
				TimestampFormat:   time.RFC3339Nano,
				DisableHTMLEscape: true,
			}})
		logger.Hooks.Add(hook)
	} else if logStore == loggerStoreLoki {
		//loki
		if lokiClient == nil {
			return logger
		}

		logger.Hooks.Add(lokiClient)
	}
	logger.Out = ioutil.Discard

	return logger
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := url.QueryUnescape(c.Request.URL.String())
		if err != nil {
			u = c.Request.URL.String()
		}

		clientInfo := GetClientInfo(c)
		log := LoggerInit().WithFields(logrus.Fields{
			RequestIDKey: GetRequestID(c),
			"Method":     c.Request.Method,
			"UserAgent":  c.Request.UserAgent(),
			"Url":        u,
			"Host":       c.Request.Host,
			"Remote":     clientInfo.IP,
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
				"UserID": GetUserID(c),
			}).Warnf("Responded %03d in %dms", c.Writer.Status(), spent)
		} else {
			log.WithFields(logrus.Fields{
				"Status": c.Writer.Status(),
				"Spent":  spent,
				"UserID": GetUserID(c),
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

// 公共Filed提取出来供不同客户端使用
func GetLoggerCommonFields() logrus.Fields {
	//host, _ := os.Hostname()
	return logrus.Fields{
		"Type":    "bird",
		"Service": "manager",
	}
}

// 初始化loki客户端， 不需要LoggerInit每次初始化loki客户端， loki客户端会开启一个协程批量发送数据
func newLokiClient(url string) (*logrusloki.Loki, error) {
	lokiClient, err := logrusloki.NewLokiDefaults(url, time.Second*5)
	if err != nil {
		return nil, fmt.Errorf("New loki logger client err:%v.", err)
	}

	fields := GetLoggerCommonFields()
	for k, v := range fields {
		lokiClient.AddData(k, v.(string))
	}

	return lokiClient, nil
}
