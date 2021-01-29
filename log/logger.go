package log

import (
	"fmt"
	"github.com/hcolde/reviewer-helper/conf"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
)

var (
	Logger    *zap.SugaredLogger
)

func getLogWriter(logFile string, maxSize, maxAge, maxBackups int) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
	})
}

// shortLine:是否显示文件及行数
func getEncoder(shortLine bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if !shortLine {
		encoderConfig.EncodeCaller = nil
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func initLog(logDir, fileName string, shortLine bool) *zap.SugaredLogger {
	if _, err := os.Stat(logDir); err != nil {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	logFile := path.Join(logDir, fileName)
	writeSyncer := getLogWriter(logFile, conf.Conf.Log.MaxSize, conf.Conf.Log.MaxAge, conf.Conf.Log.MaxBackups)
	encoder := getEncoder(shortLine)

	level := zapcore.DebugLevel
	switch conf.Conf.Log.LogLevel {
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.DebugLevel
	}
	core := zapcore.NewCore(encoder, writeSyncer, level)

	return zap.New(core, zap.AddCaller()).Sugar()
}

func init() {
	Logger = initLog(conf.Conf.Log.LogDir, conf.Conf.Log.LogFileName, true)
}
