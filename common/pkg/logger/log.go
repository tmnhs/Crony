package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/tmnhs/crony/common/pkg/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"

	"os"
)

//默认logger处理器
var _defaultLogger *zap.Logger

func Init(projectName string, level string, format, prefix, director string, showLine bool, encodeLevel string, stacktraceKey string, logInConsole bool) (logger *zap.Logger) {
	if ok := utils.Exists(fmt.Sprintf("%s/%s", projectName, director)); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", director)
		_ = os.Mkdir(fmt.Sprintf("%s/%s", projectName, director), os.ModePerm)
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
	cores := make([]zapcore.Core, 0)
	switch level {
	case "info":
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_info.log", projectName, director), infoPriority))
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_warn.log", projectName, director), warnPriority))
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_error.log", projectName, director), errorPriority))
	case "warn":
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_warn.log", projectName, director), warnPriority))
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_error.log", projectName, director), errorPriority))
	case "error":
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_error.log", projectName, director), errorPriority))
	default:
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_debug.log", projectName, director), debugPriority))
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_info.log", projectName, director), infoPriority))
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_warn.log", projectName, director), warnPriority))
		cores = append(cores, getEncoderCore(logInConsole, prefix, format, encodeLevel, stacktraceKey, fmt.Sprintf("%s/%s/server_error.log", projectName, director), errorPriority))
	}
	logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())

	if showLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	_defaultLogger = logger
	return logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig(prefix, encodeLevel, stacktraceKey string) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: stacktraceKey,
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(prefix + utils.TimeFormatDateV4))
		},
		//EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case encodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case encodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case encodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case encodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder(prefix, format, encodeLevel, stacktraceKey string) zapcore.Encoder {
	if format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(prefix, encodeLevel, stacktraceKey))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(prefix, encodeLevel, stacktraceKey))
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(logInConsole bool, prefix, format, encodeLevel, stacktraceKey string, fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := getWriteSyncer(logInConsole, fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(prefix, format, encodeLevel, stacktraceKey), writer, level)
}

func getWriteSyncer(logInConsole bool, file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, // 日志文件的位置
		MaxSize:    10,   // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 200,  // 保留旧文件的最大个数
		MaxAge:     30,   // 保留旧文件的最大天数
		Compress:   true, // 是否压缩/归档旧文件
	}

	if logInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	}
	return zapcore.AddSync(lumberJackLogger)
}

/*func Infof(msg string, fields ...interface{}) {
	if len(fields) == 0 {
		_defaultLogger.Info(msg)
		return
	}
	_defaultLogger.Info(fmt.Sprintf(msg, fields))
}

func Info(msg string, fields ...zap.Field) {
	_defaultLogger.Info(msg, fields...)
}

func Debugf(msg string, fields ...interface{}) {
	if len(fields) == 0 {
		_defaultLogger.Debug(msg)
		return
	}
	_defaultLogger.Debug(fmt.Sprintf(msg, fields))
}

func Debug(msg string, fields ...zap.Field) {
	_defaultLogger.Debug(msg, fields...)
}

func Warnf(msg string, fields ...interface{}) {
	if len(fields) == 0 {
		_defaultLogger.Warn(msg)
		return
	}
	_defaultLogger.Warn(fmt.Sprintf(msg, fields))
}

func Warn(msg string, fields ...zap.Field) {
	_defaultLogger.Warn(msg, fields...)
}

func Errorf(msg string, fields ...interface{}) {
	if len(fields) == 0 {
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
	if len(fields) == 0 {
		_defaultLogger.Fatal(msg)
		return
	}
	_defaultLogger.Fatal(fmt.Sprintf(msg, fields))
}*/

func Sync() error {
	return _defaultLogger.Sync()
}

//Shutdown 全局关闭
func Shutdown() {
	_defaultLogger.Sync()
}
func GetLogger() *zap.Logger {
	return _defaultLogger
}
