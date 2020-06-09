package testZapLogger_test

import (
	"testing"
	"time"

	"go.uber.org/zap/zapcore"
	"golang.org/x/xerrors"

	"github.com/goroumaru/test-code/testZapLogger"
	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	// デフォルト設定でロギングする
	testZapLogger.Logger()
}

func TestLogger2(t *testing.T) {

	// ユーザ設定でロギングする（ログ出力および標準出力の２箇所へ）
	logger, logFile := testZapLogger.NewLogger()
	defer logFile.Close()

	// DEBUG LEVEL
	logger.Debug("Please write messages here!") // msgキーへ登録される

	// INFO LEVEL ,log & stdout
	logger.Info("Please write messages here!", zap.Stack("StackTrace"))     // このInfoのスタックトレースを出力する
	logger.Info("Please write messages here!", zap.String("key", "value"))  // "key":"value"形式
	logger.Warn("Please write messages here!", zap.Time("now", time.Now())) // 時刻追加

	// WARN LEVEL ,log & stdout
	logger.Warn("Please write messages here!")

	// ERROR LEVEL ,log & stdout
	logger.Error("Please write messages here!", zap.Error(doEverything()))  // xerrorと組み合わせて使うと、スタックトレース内容も出力できる
	logger.Error("Please write messages here!", zap.Error(doEverything2())) // xerrorでwrap使用しても、トレースできる

	// ユーザー設定したログ内容
	user := &user{
		ID:  1,
		Msg: "testです",
	}

	// user objectを挿入する、クロージャ
	logger.Warn("Please write messages here!", zap.Object("userObj", zapcore.ObjectMarshalerFunc(func(inner zapcore.ObjectEncoder) error {
		inner.AddInt("id", user.ID)
		inner.AddString("msg", user.Msg)
		return nil
	})))

	// user objectを挿入する、レシーバー
	logger.Error("Please write messages here!", zap.Object("userObj", user))
}

type user struct {
	ID  int
	Msg string
}

// zap.Objectの引数
func (u *user) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", u.ID)
	enc.AddString("msg", u.Msg)
	return nil
}

// xerrorsでエラースタックする
func doSomething() error {
	return xerrors.New("error occured here!")
}
func doAnything() error {
	err := doSomething()
	return xerrors.Errorf("Second: %v", err)
}
func doEverything() error {
	err := doAnything()
	return xerrors.Errorf("Third: %v", err)
}

// xerrorsでラップする `: %w`
func doSomething2() error {
	return xerrors.New("error occured here!")
}
func doAnything2() error {
	err := doSomething2()
	return xerrors.Errorf("Second: %w", err)
}
func doEverything2() error {
	err := doAnything2()
	return xerrors.Errorf("Third: %w", err)
}
