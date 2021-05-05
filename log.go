package log

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	defaultLogger = NewLogger(DebugLevel, "default")
)

// Level type
type Level uint32

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
	timeFormat             = "2006-01-02 15:04:05.000"
)

var (
	// Used for caller information initialisation
	callerInitOnce     sync.Once
	logrusPackage      string
	minimumCallerDepth = 1
	loggers            = make(map[string]*MyLogger)
)

type Fields map[string]interface{}

func (fields Fields) String() string {
	str := make([]string, 0)

	for k, v := range fields {
		str = append(str, fmt.Sprintf("%s=%+v", k, v))
	}

	return strings.Join(str, " ")
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}

func getFuncName(f string) string {
	n := strings.LastIndex(f, "/")
	if n == -1 {
		return f
	}
	return f[n+1:]
}

// getCaller retrieves the name of the first non-logrus calling function
func getCaller() *runtime.Frame {
	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		pcs := make([]uintptr, maximumCallerDepth)
		_ = runtime.Callers(0, pcs)

		// dynamic get the package name and the minimum caller depth
		for i := 0; i < maximumCallerDepth; i++ {
			funcName := runtime.FuncForPC(pcs[i]).Name()
			if strings.Contains(funcName, "getCaller") {
				logrusPackage = getPackageName(funcName)
				break
			}
		}

		minimumCallerDepth = knownLogrusFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		// find the function which is not logrus and ion-log
		if !strings.Contains(pkg, "logrus") && pkg != logrusPackage {
			return &f //nolint:scopelint
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// Infof logs a formatted info level log to the console
func Infof(format string, v ...interface{}) { defaultLogger.Infof(format, v...) }

// Tracef logs a formatted debug level log to the console
func Tracef(format string, v ...interface{}) { defaultLogger.Tracef(format, v...) }

// Debugf logs a formatted debug level log to the console
func Debugf(format string, v ...interface{}) { defaultLogger.Debugf(format, v...) }

// Warnf logs a formatted warn level log to the console
func Warnf(format string, v ...interface{}) { defaultLogger.Warnf(format, v...) }

// Errorf logs a formatted error level log to the console
func Errorf(format string, v ...interface{}) { defaultLogger.Errorf(format, v...) }

// Panicf logs a formatted panic level log to the console.
// The panic() function is called, which stops the ordinary flow of a goroutine.
func Panicf(format string, v ...interface{}) { defaultLogger.Panicf(format, v...) }

func Init(level string) {
	l := logrus.DebugLevel
	switch level {
	case "trace":
		l = logrus.TraceLevel
	case "debug":
		l = logrus.DebugLevel
	case "info":
		l = logrus.InfoLevel
	case "warn":
		l = logrus.WarnLevel
	case "error":
		l = logrus.ErrorLevel
	}
	defaultLogger.SetLevel(l)
}

// get goroutine id
// func getGID() uint64 {
// b := make([]byte, 64)
// b = b[:runtime.Stack(b, false)]
// b = bytes.TrimPrefix(b, []byte("goroutine "))
// b = b[:bytes.IndexByte(b, ' ')]
// n, _ := strconv.ParseUint(string(b), 10, 64)
// return n
// }

type MyLogger struct {
	logger *logrus.Logger
	level  Level
	prefix string
}

func (ml *MyLogger) Level() string {
	switch ml.level {
	case PanicLevel:
		return "Panic"
	case FatalLevel:
		return "Fatal"
	case ErrorLevel:
		return "Error"
	case WarnLevel:
		return "Warn"
	case InfoLevel:
		return "Info"
	case DebugLevel:
		return "Debug"
	case TraceLevel:
		return "Trace"
	}
	return "Unkown"
}

func (ml *MyLogger) Prefix() string {
	return ml.prefix
}

func (ml *MyLogger) SetLevel(level Level) {
	ml.logger.SetLevel(logrus.Level(level))
}

func NewLogger(level Level, prefix string) *logrus.Logger {
	if logger, found := loggers[prefix]; found {
		return logger.logger
	}
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetReportCaller(true)
	l.SetLevel(logrus.Level(level))
	l.SetFormatter(&TextFormatter{
		Prefix:          prefix,
		FullTimestamp:   true,
		TimestampFormat: timeFormat,
	})

	loggers[prefix] = &MyLogger{
		logger: l,
		level:  level,
		prefix: prefix,
	}
	return l
}

func NewLoggerWithFields(level Level, prefix string, fields Fields) *logrus.Logger {
	if logger, found := loggers[prefix]; found {
		return logger.logger
	}
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetReportCaller(true)
	l.SetLevel(logrus.Level(level))
	l.SetFormatter(&TextFormatter{
		Prefix:          prefix,
		Fields:          fields,
		FullTimestamp:   true,
		TimestampFormat: timeFormat,
	})

	loggers[prefix] = &MyLogger{
		logger: l,
		level:  level,
		prefix: prefix,
	}

	return l
}

func SetLogLevel(prefix string, level Level) error {
	if l, found := loggers[prefix]; found {
		l.level = level
		l.logger.SetLevel(logrus.Level(level))
		return nil
	}
	return fmt.Errorf("logger [%v] not found", prefix)
}

func GetLoggers() map[string]*MyLogger {
	return loggers
}
