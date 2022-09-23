package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/tmnhs/crony/common/models"
	"github.com/tmnhs/crony/common/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"os"
)

//默认logger处理器
var _defaultLogger *zap.Logger

func Init(l *models.Log, projectName string) (logger *zap.Logger) {
	if ok := utils.Exists(fmt.Sprintf("%s/%s",projectName,l.Director)); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", l.Director)
		_ = os.Mkdir(fmt.Sprintf("%s/%s",projectName,l.Director), os.ModePerm)
	}
	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	cores := [...]zapcore.Core{
		getEncoderCore(l, fmt.Sprintf("%s/%s/server_debug.log",projectName, l.Director), debugPriority),
		getEncoderCore(l, fmt.Sprintf("%s/%s/server_info.log",projectName, l.Director), infoPriority),
		getEncoderCore(l, fmt.Sprintf("%s/%s/server_warn.log",projectName, l.Director), warnPriority),
		getEncoderCore(l, fmt.Sprintf("%s/%s/server_error.log",projectName, l.Director), errorPriority),
	}
	logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())

	if l.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	_defaultLogger = logger
	return logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig(l *models.Log) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  l.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case l.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case l.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case l.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case l.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder(l *models.Log) zapcore.Encoder {
	if l.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(l))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(l))
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(l *models.Log, fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := GetWriteSyncer(l, fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(l), writer, level)
}

func GetWriteSyncer(l *models.Log, file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, // 日志文件的位置
		MaxSize:    10,   // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 200,  // 保留旧文件的最大个数
		MaxAge:     30,   // 保留旧文件的最大天数
		Compress:   true, // 是否压缩/归档旧文件
	}

	if l.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}

func Infof(msg string, fields ...interface{}) {
	if len(fields)==0{
		_defaultLogger.Info(msg)
		return
	}
	_defaultLogger.Info(fmt.Sprintf(msg, fields))
}

func Info(msg string, fields ...zap.Field) {
	_defaultLogger.Info(msg, fields...)
}

func Debugf(msg string, fields ...interface{}) {
	if len(fields)==0{
		_defaultLogger.Debug(msg)
		return
	}
	_defaultLogger.Debug(fmt.Sprintf(msg, fields))
}

func Debug(msg string, fields ...zap.Field) {
	_defaultLogger.Debug(msg, fields...)
}

func Warnf(msg string, fields ...interface{}) {
	if len(fields)==0 {
		_defaultLogger.Warn(msg)
		return
	}
	_defaultLogger.Warn(fmt.Sprintf(msg, fields))
}

func Warn(msg string, fields ...zap.Field) {
	_defaultLogger.Warn(msg, fields...)
}

func Errorf(msg string, fields ...interface{}) {
	if len(fields)==0 {
		_defaultLogger.Error(msg)
		return
	}
	_defaultLogger.Error(fmt.Sprintf(msg, fields))
}

func Error(msg string, fields ...zap.Field) {
	_defaultLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	_defaultLogger.Fatal(msg, fields...)
}

func Fatalf(msg string, fields ...interface{}) {
	if len(fields)==0 {
		_defaultLogger.Fatal(msg)
		return
	}
	_defaultLogger.Fatal(fmt.Sprintf(msg, fields))
}

func Sync() error {
	return _defaultLogger.Sync()
}

//Shutdown 全局关闭
func Shutdown() {
	_defaultLogger.Sync()
}
