package frameworks

import (
	"log"
	"strings"

	"github.com/riomhaire/lightauth2/usecases"
)

type ConsoleLogger struct {
	loggingLevel int
}

func NewConsoleLogger(level string) ConsoleLogger {
	l := strings.ToLower(level)
	v := usecases.Debug
	switch l {
	case "trace":
		v = usecases.Trace
	case "debug":
		v = usecases.Debug
	case "info":
		v = usecases.Info
	case "warn":
		v = usecases.Warn
	case "error":
		v = usecases.Error
	}
	return ConsoleLogger{v}
}

func toLoggingLevel(l int) string {
	switch l {
	case usecases.Trace:
		return "trace"
	case usecases.Debug:
		return "debug"
	case usecases.Info:
		return "info"
	case usecases.Warn:
		return "warn"
	case usecases.Error:
		return "error"
	}
	return "unknown"

}

func (d ConsoleLogger) Log(level int, message string) {
	if level >= d.loggingLevel {
		log.Printf("[%s] %s\n", toLoggingLevel(level), message)
	}
}
