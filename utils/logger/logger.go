package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger :
type Logger struct {
	log     *zap.Logger
	logFile *os.File
	mux     sync.Mutex
}

// NewLogger :
// https://k1low.hatenablog.com/entry/2018/08/15/100000
func NewLogger(logFilePath string, logLevel string) *Logger {
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
	lf, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		fmt.Printf("Failed to open log file: %v", err)
	}

	// コンソール・コア生成
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // コンソール出力
		zapcore.AddSync(os.Stdout),               // io.Writerをzapcore.WriteSyncerに変換
		setLevel(logLevel),                       // 出力レベル設定
	)

	// ログ・コア生成
	logCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 構造化ログ（JSON）
		zapcore.AddSync(lf),
		setLevel(logLevel),
	)

	// コア結合
	log := zap.New(zapcore.NewTee(
		consoleCore,
		logCore,
	))
	return &Logger{log, lf, sync.Mutex{}} // 呼び出し元でdefer Close()メソッドし、ファイルを閉じる
}

// Close : close logger
func (l *Logger) Close() {
	l.Info("Close logger.", nil)

	l.mux.Lock()
	defer l.mux.Unlock()
	if err := l.syncFile(); err != nil {
		l.Info("Log file can not be synchronized.", zap.Error(err))
	}
	if err := l.closeFile(); err != nil {
		l.Info("Log file can not be closed.", zap.Error(err))
	}
}

// Debug : use zap logger. If val don't exist, set "nil".
func (l *Logger) Debug(message string, val interface{}, keys ...string) {
	assertFields(l.log.Debug, message, val, keys...)
}

// Info : use zap logger. If val don't exist, set "nil".
func (l *Logger) Info(message string, val interface{}, keys ...string) {
	assertFields(l.log.Info, message, val, keys...)
}

// Warn : use zap logger. If val don't exist, set "nil".
func (l *Logger) Warn(message string, val interface{}, keys ...string) {
	assertFields(l.log.Warn, message, val, keys...)
}

// Error : use zap logger. If val don't exist, set "nil".
func (l *Logger) Error(message string, val interface{}, keys ...string) {
	assertFields(l.log.Error, message, val, keys...)
}

// https://qiita.com/tsurumiii/items/0294feebc0216b185765
func assertFields(log func(msg string, fields ...zapcore.Field), message string, val interface{}, keys ...string) {
	var key string
	if len(keys) != 0 {
		key = keys[0]
	}
	log(message, zap.Any(key, val))
}

func setLevel(logLevel string) (level zapcore.Level) {
	switch logLevel {
	case "Debug":
		level = zapcore.DebugLevel
	case "Info":
		level = zapcore.InfoLevel
	case "Warn":
		level = zapcore.WarnLevel
	case "Error":
		level = zapcore.ErrorLevel
	}
	return
}

// Flushing the file system's in-memory copy of recently written data to disk.
func (l *Logger) syncFile() error {
	return l.logFile.Sync()
}

// Close log file
func (l *Logger) closeFile() error {
	if l.logFile == nil {
		return nil
	}
	err := l.logFile.Close()
	l.logFile = nil // closeしたので、nilとする
	return err
}
