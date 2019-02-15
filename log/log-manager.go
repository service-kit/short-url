package log

import (
	"github.com/service-kit/short-url/common"
	"github.com/service-kit/short-url/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
)

type LogManager struct {
	logger *zap.Logger
}

var m *LogManager
var once sync.Once

func GetInstance() *LogManager {
	once.Do(func() {
		m = &LogManager{}
	})
	return m
}

func (self *LogManager) InitManager() (err error) {
	logPath, _ := config.GetInstance().GetConfig("LOG_FILE_PATH")
	logLevel, _ := config.GetInstance().GetConfig("LOG_LEVEL")
	if "" == logPath {
		cfg := zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		self.logger, err = cfg.Build()
		if nil != err {
			return
		}
	} else {
		hook := lumberjack.Logger{
			Filename:   logPath, // 日志文件路径
			MaxSize:    1024,    // megabytes
			MaxBackups: 100,     // 最多保留3个备份
			MaxAge:     30,      //days
			Compress:   true,    // 是否压缩 disabled by default
		}
		writeSyncer := zapcore.AddSync(&hook)
		var level zapcore.Level
		switch logLevel {
		case common.LL_DEBUG:
			level = zap.DebugLevel
		case common.LL_INFO:
			level = zap.InfoLevel
		case common.LL_ERROR:
			level = zap.ErrorLevel
		default:
			level = zap.InfoLevel
		}
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			writeSyncer,
			level,
		)
		self.logger = zap.New(core)
	}
	self.logger = self.logger.Named("asr-auth")
	return
}

func (self *LogManager) FinishProcess() {
	if nil != self.logger {
		self.logger.Sync()
	}
}

func (self *LogManager) GetLogger() *zap.Logger {
	return self.logger
}
