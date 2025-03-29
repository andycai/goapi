package log

// 对 zap 进行封装

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// Config 日志配置
type Config struct {
	Level      string // 日志级别
	Filename   string // 日志文件路径
	MaxSize    int    // 每个日志文件最大尺寸（MB）
	MaxBackups int    // 保留的旧日志文件最大数量
	MaxAge     int    // 保留的旧日志文件最大天数
	Compress   bool   // 是否压缩旧日志文件
}

// InitLogger 初始化日志配置
func InitLogger(cfg *Config) {
	if cfg == nil {
		Logger, _ = zap.NewProduction()
		Sugar = Logger.Sugar()
		return
	}

	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 配置输出
	var core zapcore.Core
	if cfg.Filename != "" {
		// 文件输出
		writer := &lumberjack.Logger{
			Filename:   cfg.Filename,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(writer),
			level,
		)
	} else {
		// 控制台输出
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	}

	// 创建Logger
	Logger = zap.New(core, zap.AddCaller())
	Sugar = Logger.Sugar()
}

func init() {
	InitLogger(nil)
}

// WithFields 添加自定义字段
func WithFields(fields map[string]interface{}) *zap.SugaredLogger {
	if len(fields) == 0 {
		return Sugar
	}
	zapFields := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		zapFields = append(zapFields, k, v)
	}
	return Sugar.With(zapFields...)
}

// Debug 输出Debug级别日志
func Debug(args ...any) {
	Sugar.Debug(args...)
}

// Debugf 输出Debug级别格式化日志
func Debugf(template string, args ...any) {
	Sugar.Debugf(template, args...)
}

// Info 输出Info级别日志
func Info(args ...any) {
	Sugar.Info(args...)
}

// Infof 输出Info级别格式化日志
func Infof(template string, args ...any) {
	Sugar.Infof(template, args...)
}

// Warn 输出Warn级别日志
func Warn(args ...any) {
	Sugar.Warn(args...)
}

// Warnf 输出Warn级别格式化日志
func Warnf(template string, args ...any) {
	Sugar.Warnf(template, args...)
}

// Error 输出Error级别日志
func Error(args ...any) {
	Sugar.Error(args...)
}

// Errorf 输出Error级别格式化日志
func Errorf(template string, args ...any) {
	Sugar.Errorf(template, args...)
}

// Fatal 输出Fatal级别日志
func Fatal(args ...any) {
	Sugar.Fatal(args...)
}

// Fatalf 输出Fatal级别格式化日志
func Fatalf(template string, args ...any) {
	Sugar.Fatalf(template, args...)
}

// Sync 同步缓存的日志
func Sync() error {
	return Logger.Sync()
}
