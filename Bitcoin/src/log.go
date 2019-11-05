package main
/*
	go日志框架zap配置
*/
import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

/*
	定义全局变量
*/
var Logger *zap.Logger
var Sugar *zap.SugaredLogger

/*
	加载
*/
func init() {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "log/im.log", // ⽇志⽂件路径
		MaxSize:    1024,         //megabytes
		MaxBackups: 3,            // 最多保留3个备份
		MaxAge:     28,           //days
		Compress:   true,
	})
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(NewProductionEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout),
			w),
		zap.DebugLevel,
	)
	Logger = zap.New(core, zap.AddCaller())
	Sugar = Logger.Sugar()
}

func NewProductionEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

/*
	设置时间格式
*/
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
