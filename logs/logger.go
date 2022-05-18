package logger

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const (
	LOGGER_TIME_FORMAT string = "2006-01-02 15:04:05.999999999 -07:00"

	// Logger LEVEL reference from logrus.Level
	PanicLevel = log.PanicLevel
	FatalLevel = log.FatalLevel
	ErrorLevel = log.ErrorLevel
	WarnLevel  = log.WarnLevel
	InfoLevel  = log.InfoLevel
	DebugLevel = log.DebugLevel

	// DMS Server components
	MQTT      ServerComponent = "MQTT"
	SQLSERVER ServerComponent = "SQLSERVER"
	DMSSERVER ServerComponent = "DMS_SERVER"
	GINROUTER ServerComponent = "GIN_ROUTER"
)

var logger = log.New()

type (
	mqttLogger interface {
		Println(v ...interface{})
		Printf(format string, v ...interface{})
	}
	NOOPLogger struct {
		prefix string
		level  log.Level
	}
	loggerEntry struct {
		entry *log.Entry
	}
	ServerComponent string
	LoggerFields    log.Fields
)

func InitLogger(logFilePath string) {
	logger.SetFormatter(&LogFormatter{
		TimestampFormat: LOGGER_TIME_FORMAT,
		WithCallerField: true,
	})
	logFile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		Fatal("Cannot open ", logFilePath)
	}
	logger.SetLevel(log.DebugLevel)
	logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func createFieldsWithComponent(comp ServerComponent, fields LoggerFields) (logFields map[string]interface{}) {
	if fields == nil || len(fields) == 0 {
		logFields = make(map[string]interface{}, 1)
		logFields["Component"] = comp
	} else {
		logFields = make(map[string]interface{}, len(fields)+1)
		logFields["Component"] = comp
		for k, v := range fields {
			logFields[k] = v
		}
	}
	return logFields
}

func (n *NOOPLogger) Println(v ...interface{}) {
	entry := &loggerEntry{
		entry: logger.WithFields(createFieldsWithComponent(MQTT, LoggerFields{
			"MqttLevel": n.prefix,
		})),
	}
	entry.log(n.level, v...)
}

func (n *NOOPLogger) Printf(format string, v ...interface{}) {
	entry := &loggerEntry{
		entry: logger.WithFields(createFieldsWithComponent(MQTT, LoggerFields{
			"MqttLevel": n.prefix,
		})),
	}
	entry.logf(n.level, format, v...)
}

func NewMqttLogger(prefix string, logLevel log.Level) *NOOPLogger {
	return &NOOPLogger{
		prefix: prefix,
		level:  log.Level(logLevel),
	}
}

func (en *loggerEntry) log(level log.Level, args ...interface{}) {
	switch level {
	case PanicLevel:
		en.entry.Panicln(args...)
	case FatalLevel:
		en.entry.Fatalln(args...)
	case ErrorLevel:
		en.entry.Errorln(args...)
	case WarnLevel:
		en.entry.Warnln(args...)
	case InfoLevel:
		en.entry.Infoln(args...)
	case DebugLevel:
		en.entry.Debugln(args...)
	}
}

func (en *loggerEntry) logf(level log.Level, format string, args ...interface{}) {
	switch level {
	case PanicLevel:
		en.entry.Panicf(format, args...)
	case FatalLevel:
		en.entry.Fatalf(format, args...)
	case ErrorLevel:
		en.entry.Errorf(format, args...)
	case WarnLevel:
		en.entry.Warnf(format, args...)
	case InfoLevel:
		en.entry.Infof(format, args...)
	case DebugLevel:
		en.entry.Debugf(format, args...)
	}
}

func Fatal(args ...interface{}) {
	logger.Fatalln(args...)
}

func Debug(args ...interface{}) {
	logger.Debugln(args...)
}

func Info(args ...interface{}) {
	logger.Infoln(args...)
}

func Error(args ...interface{}) {
	logger.Errorln(args...)
}

func Warn(args ...interface{}) {
	logger.Warnln(args...)
}

func Panic(args ...interface{}) {
	logger.Panicln(args...)
}

func Trace(args ...interface{}) {
	logger.Traceln(args...)
}

func GinLogger(notLogged ...string) gin.HandlerFunc {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, p := range notLogged {
			skip[p] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()
		c.Next()
		stop := time.Since(start)
		latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000000.0))
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		referer := c.Request.Referer()
		dataLength := c.Writer.Size()
		if dataLength < 0 {
			dataLength = 0
		}

		if _, ok := skip[path]; ok {
			return
		}

		ginLogFields := createFieldsWithComponent(GINROUTER, LoggerFields{
			"hostname":   hostname,
			"statusCode": statusCode,
			"latency":    latency,
			"clientIP":   clientIP,
			"method":     c.Request.Method,
			"path":       path,
			"referer":    referer,
			"dataLength": dataLength,
			"userAgent":  userAgent,
		})
		entry := logger.WithFields(ginLogFields)

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			msg := fmt.Sprintf("%s - %s \"%s %s\" %d %d \"%s\" \"%s\" (%dms)", clientIP, hostname, c.Request.Method, path, statusCode, dataLength, referer, userAgent, latency)
			if statusCode >= http.StatusInternalServerError {
				entry.Error(msg)
			} else if statusCode >= http.StatusBadRequest {
				entry.Warn(msg)
			} else {
				entry.Info(msg)
			}
		}
	}
}

func LogWithFields(comp ServerComponent, level log.Level, fields LoggerFields, message ...interface{}) {
	logFields := createFieldsWithComponent(comp, fields)
	entry := &loggerEntry{
		entry: logger.WithFields(logFields),
	}
	entry.log(level, message...)
}

func LogfWithFields(comp ServerComponent, level log.Level, fields LoggerFields, messageFormat string, args ...interface{}) {
	logFields := createFieldsWithComponent(comp, fields)
	entry := &loggerEntry{
		entry: logger.WithFields(logFields),
	}
	entry.logf(level, messageFormat, args...)
}

func LogWithoutFields(comp ServerComponent, level log.Level, message ...interface{}) {
	entry := &loggerEntry{
		entry: logger.WithFields(createFieldsWithComponent(comp, nil)),
	}
	entry.log(level, message...)
}

func LogfWithoutFields(comp ServerComponent, level log.Level, messageFormat string, args ...interface{}) {
	entry := &loggerEntry{
		entry: logger.WithFields(createFieldsWithComponent(comp, nil)),
	}
	entry.logf(level, messageFormat, args...)
}
