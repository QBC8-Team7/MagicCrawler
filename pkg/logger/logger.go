package logger

import (
	"log"
	"os"

	"github.com/QBC8-Team7/MagicCrawler/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger methods interface
type Logger interface {
	InitLogger(appLogFilePath, sysLogFilePath string)
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	LogSystemInfo(args ...interface{})
	LogSystemInfof(template string, args ...interface{})
}

// Logger
type AppLogger struct {
	cfg         *config.Config
	sugarLogger *zap.SugaredLogger
	sysLogger   *zap.SugaredLogger // Separate logger for system metrics
}

// App Logger constructor
func NewAppLogger(cfg *config.Config) *AppLogger {
	return &AppLogger{cfg: cfg}
}

// Mapping config logger levels to zap levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *AppLogger) getLoggerLevel(cfg *config.Config) zapcore.Level {
	level, exists := loggerLevelMap[cfg.Logger.Level]
	if !exists {
		return zapcore.DebugLevel
	}
	return level
}

// InitLogger initializes the main and system loggers
func (l *AppLogger) InitLogger(appLogFilePath, sysLogFilePath string) {
	logLevel := l.getLoggerLevel(l.cfg)

	appFile, err := os.OpenFile(appLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cannot open application log file:", err)
	}
	appLogWriter := zapcore.AddSync(appFile)

	sysFile, err := os.OpenFile(sysLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Cannot open system log file:", err)
	}
	sysLogWriter := zapcore.AddSync(sysFile)

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	appCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		appLogWriter,
		zap.NewAtomicLevelAt(logLevel),
	)

	sysCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		sysLogWriter,
		zap.NewAtomicLevelAt(zapcore.InfoLevel),
	)

	l.sugarLogger = zap.New(appCore, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	l.sysLogger = zap.New(sysCore, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()

	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error("Failed to sync main logger:", err)
	}
	if err := l.sysLogger.Sync(); err != nil {
		l.sysLogger.Error("Failed to sync system logger:", err)
	}
}

func (l *AppLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *AppLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *AppLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *AppLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *AppLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *AppLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *AppLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *AppLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *AppLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *AppLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *AppLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *AppLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

func (l *AppLogger) LogSystemInfo(args ...interface{}) {
	l.sysLogger.Info(args...)
}

func (l *AppLogger) LogSystemInfof(template string, args ...interface{}) {
	l.sysLogger.Infow(template, args...)
}
