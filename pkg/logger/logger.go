package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap logger
type Logger struct {
	*zap.Logger
}

// Config represents logger configuration
type Config struct {
	Level  string
	Format string
	Output string
}

// New creates a new logger
func New(cfg Config) (*Logger, error) {
	// Parse log level
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	// Configure encoder
	var encoderConfig zapcore.EncoderConfig
	if cfg.Format == "json" {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create encoder
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure output
	var writer zapcore.WriteSyncer
	if cfg.Output == "" || cfg.Output == "stdout" {
		writer = zapcore.AddSync(os.Stdout)
	} else if cfg.Output == "stderr" {
		writer = zapcore.AddSync(os.Stderr)
	} else {
		file, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		writer = zapcore.AddSync(file)
	}

	// Create core
	core := zapcore.NewCore(encoder, writer, level)

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{Logger: logger}, nil
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{Logger: l.With(fields...)}
}

// WithComponent adds a component field
func (l *Logger) WithComponent(component string) *Logger {
	return l.WithFields(zap.String("component", component))
}

// WithConnectionID adds a connection ID field
func (l *Logger) WithConnectionID(id uint32) *Logger {
	return l.WithFields(zap.Uint32("conn_id", id))
}

// WithError adds an error field
func (l *Logger) WithError(err error) *Logger {
	return l.WithFields(zap.Error(err))
}

// LogConnectionEvent logs a connection event
func (l *Logger) LogConnectionEvent(event string, connID uint32, remoteAddr string) {
	l.Info(event,
		zap.String("event", event),
		zap.Uint32("conn_id", connID),
		zap.String("remote_addr", remoteAddr),
	)
}

// LogDataTransfer logs data transfer statistics
func (l *Logger) LogDataTransfer(connID uint32, bytesSent, bytesReceived uint64) {
	l.Info("data_transfer",
		zap.Uint32("conn_id", connID),
		zap.Uint64("bytes_sent", bytesSent),
		zap.Uint64("bytes_received", bytesReceived),
	)
}

// LogHandshake logs a handshake event
func (l *Logger) LogHandshake(success bool, remoteAddr string, reason string) {
	if success {
		l.Info("handshake_success",
			zap.String("remote_addr", remoteAddr),
		)
	} else {
		l.Warn("handshake_failed",
			zap.String("remote_addr", remoteAddr),
			zap.String("reason", reason),
		)
	}
}

// LogPacket logs packet information
func (l *Logger) LogPacket(direction string, packetType string, size int) {
	l.Debug("packet",
		zap.String("direction", direction),
		zap.String("type", packetType),
		zap.Int("size", size),
	)
}

// Default logger instance
var defaultLogger *Logger

// InitDefault initializes the default logger
func InitDefault(cfg Config) error {
	logger, err := New(cfg)
	if err != nil {
		return err
	}
	defaultLogger = logger
	return nil
}

// Default returns the default logger
func Default() *Logger {
	if defaultLogger == nil {
		// Create a basic logger if not initialized
		logger, _ := New(Config{
			Level:  "info",
			Format: "console",
			Output: "stdout",
		})
		defaultLogger = logger
	}
	return defaultLogger
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	Default().Info(msg, fields...)
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	Default().Debug(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	Default().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	Default().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	Default().Fatal(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	if defaultLogger != nil {
		return defaultLogger.Sync()
	}
	return nil
}
