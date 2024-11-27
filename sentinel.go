package sentinel

import (
	"context"
	"log/slog"
	"os"
	"time"
)

//go:generate ggenums
//enum:name=OutputPath values=stdout,stderr

type Level = slog.Level

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type Logger struct {
	*slog.Logger
	defaultFields []slog.Attr
}

type Config struct {
	Level       Level
	TimeFormat  string
	JSONOutput  bool
	OutputPath  string
	ServiceName string
}

func DefaultConfig() Config {
	return Config{
		Level:       LevelInfo,
		TimeFormat: time.RFC3339,
		JSONOutput:  false,
		OutputPath:  OutputPathStdout.String(),
		ServiceName: "service",
	}
}

func New(cfg Config) (*Logger, error) {
	var handler slog.Handler

	var output *os.File
	switch cfg.OutputPath {
	case OutputPathStdout.String():
		output = os.Stdout
	case OutputPathStderr.String():
		output = os.Stderr
	default:
		var err error
		output, err = os.OpenFile(cfg.OutputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
	}

	opts := &slog.HandlerOptions{
		Level: cfg.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				return slog.String(a.Key, t.Format(cfg.TimeFormat))
			}
			return a
		},
	}

	// Create appropriate handler based on output format
	if cfg.JSONOutput {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	// Create logger with default fields
	defaultFields := []slog.Attr{
		slog.String("service", cfg.ServiceName),
	}

	return &Logger{
		Logger:        slog.New(handler),
		defaultFields: defaultFields,
	}, nil
}

// With creates a new Logger with additional fields
func (l *Logger) With(fields ...any) *Logger {
	return &Logger{
		Logger:        l.Logger.With(fields...),
		defaultFields: l.defaultFields,
	}
}

// WithError creates a new Logger with an error field
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}
	return l.With("error", err)
}

// WithContext creates a new Logger with context values
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Example: extract trace ID from context
	if traceID := ctx.Value("trace_id"); traceID != nil {
		return l.With("trace_id", traceID)
	}
	return l
}

func (l *Logger) attrsToArgs(args []any) []any {
	result := make([]any, len(args)+len(l.defaultFields)*2)
	copy(result, args)
	
	offset := len(args)
	for _, attr := range l.defaultFields {
		result[offset] = attr.Key
		result[offset+1] = attr.Value
		offset += 2
	}
	
	return result
}

// Debug logs at debug level
func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, l.attrsToArgs(args)...)
}

// Info logs at info level
func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, l.attrsToArgs(args)...)
}

// Warn logs at warn level
func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, l.attrsToArgs(args)...)
}

// Error logs at error level
func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, l.attrsToArgs(args)...)
}