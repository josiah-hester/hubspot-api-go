package client

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	loggers []*log.Logger
	enabled bool
}

// NewLogger creates a new \[Logger]. The outputs variable sets the destinations to which log data will be written. The prefix appears at the beginning of each generated log line, or after the log header if the \[Lmsgprefix] flag is provided. The flag argument defines the logging properties.
func NewLogger(enabled bool, prefix string, flag int, outputs ...io.Writer) *Logger {
	if !enabled {
		return &Logger{loggers: nil, enabled: false}
	}
	if len(outputs) > 0 {
		loggers := make([]*log.Logger, len(outputs))
		for i, output := range outputs {
			loggers[i] = log.New(output, prefix, flag)
		}
		return &Logger{loggers: loggers}
	} else {
		loggers := make([]*log.Logger, 1)
		loggers[0] = log.New(os.Stdout, prefix, flag)
		return &Logger{loggers: loggers}
	}
}

func (l *Logger) Printf(format string, v ...any) {
	if l.enabled {
		for _, logger := range l.loggers {
			logger.Printf(format, v...)
		}
	}
}

func (l *Logger) Println(v ...any) {
	if l.enabled {
		for _, logger := range l.loggers {
			logger.Println(v...)
		}
	}
}

func (l *Logger) Print(v ...any) {
	if l.enabled {
		for _, logger := range l.loggers {
			logger.Print(v...)
		}
	}
}

func (l *Logger) Fatal(v ...any) {
	if l.enabled {
		for _, logger := range l.loggers {
			logger.Fatal(v...)
		}
	}
}

func (l *Logger) Fatalf(format string, v ...any) {
	if l.enabled {
		for _, logger := range l.loggers {
			logger.Fatalf(format, v...)
		}
	}
}
