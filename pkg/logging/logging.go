package logging

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/crc-org/admin-helper/pkg/constants"
)

type Logger struct {
	*slog.Logger
}

type Modification struct {
	Operation string
	IP        string
	Hosts     []string
	Caller    string
	Error     error
}

var (
	logger *Logger
	once   sync.Once
)

func GetLogger() *Logger {
	once.Do(func() {
		logger = &Logger{slog.New(slog.NewJSONHandler(logWriter(), &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))}
	})
	return logger
}

func logWriter() io.Writer {
	if path := os.Getenv(constants.LogFileEnvVar); path != "" {
		f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
		if err == nil {
			return f
		}
	}
	return os.Stdout
}

func (l *Logger) LogModification(mod Modification) {
	attrs := []any{
		"operation", mod.Operation,
		"caller", mod.Caller,
	}
	if mod.IP != "" {
		attrs = append(attrs, "ip", mod.IP)
	}
	if len(mod.Hosts) > 0 {
		attrs = append(attrs, "hosts", mod.Hosts)
	}
	if mod.Error != nil {
		attrs = append(attrs, "error", mod.Error.Error())
	}
	l.Info("hosts modification", attrs...)
}
