package testZapLogger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logPath = "./example.log"
	// configPath = "./example.json"
)

func Logger() {
	logger, _ := zap.NewDevelopment()
	logger.Info("Hello zap", zap.String("key", "value"), zap.Time("now", time.Now()))
	logger.Debug("Hello zap", zap.String("key", "value"), zap.Time("now", time.Now()))
}

// NewLogger :
// https://k1low.hatenablog.com/entry/2018/08/15/100000
func NewLogger() (*zap.Logger, *os.File) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// ログ出力ファイル
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		fmt.Printf("Failed to open log file: %v", err)
	}

	// コンソール・コア生成
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // コンソール出力
		zapcore.AddSync(os.Stdout),               // io.Writerをzapcore.WriteSyncerに変換
		zapcore.DebugLevel,                       // 出力レベル設定
	)

	// ログ・コア生成
	logCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 構造化ログ（JSON）
		zapcore.AddSync(lf),
		zapcore.DebugLevel, // 出力レベル設定
	)

	logger := zap.New(zapcore.NewTee(
		consoleCore,
		logCore,
	))

	return logger, lf // file openしているので、呼び出し元でdefer lf.Close()すること
}
