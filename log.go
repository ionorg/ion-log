package log

import (
	"bytes"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// Used for caller information initialisation
	callerInitOnce     sync.Once
	logrusPackage      string
	minimumCallerDepth = 1
)

const (
	maximumCallerDepth int = 25
	knownLogrusFrames  int = 4
)

const timeFormat = "2006-01-02 15:04:05.000"

// Config defines parameters for the logger
type Config struct {
	Level string `mapstructure:"level"`
}

// Infof logs a formatted info level log to the console
func Infof(format string, v ...interface{}) { log.Infof(format, v...) }

// Tracef logs a formatted debug level log to the console
func Tracef(format string, v ...interface{}) { log.Tracef(format, v...) }

// Debugf logs a formatted debug level log to the console
func Debugf(format string, v ...interface{}) { log.Debugf(format, v...) }

// Warnf logs a formatted warn level log to the console
func Warnf(format string, v ...interface{}) { log.Warnf(format, v...) }

// Errorf logs a formatted error level log to the console
func Errorf(format string, v ...interface{}) { log.Errorf(format, v...) }

// Panicf logs a formatted panic level log to the console.
// The panic() function is called, which stops the ordinary flow of a goroutine.
func Panicf(format string, v ...interface{}) { log.Panicf(format, v...) }

type LogFormatter struct{}

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

func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	var file string
	var len int
	// use custom getCaller because default getCaller worked bad after we wrapper it
	entry.Caller = getCaller()
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	}

	msg := fmt.Sprintf("[%s][%s:%d][%s][%s] => %s\n", timestamp, file, len, strings.ToUpper(entry.Level.String()), getFuncName(entry.Caller.Function), entry.Message)
	return []byte(msg), nil
}

// get goroutine id
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// func Init(level string) *log.Logger {
func Init(level string) {
	l := log.DebugLevel
	switch level {
	case "trace":
		l = log.TraceLevel
	case "debug":
		l = log.DebugLevel
	case "info":
		l = log.InfoLevel
	case "warn":
		l = log.WarnLevel
	case "error":
		l = log.ErrorLevel
	}
	// log := log.New()
	log.SetLevel(l)
	log.SetReportCaller(true)
	log.SetFormatter(new(LogFormatter))
	// return logger
}
