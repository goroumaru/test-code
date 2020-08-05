package logger_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/goroumaru/bot-bitbank/utils/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/xerrors"
	"gopkg.in/ini.v1"
)

func getConfig() (logPath, logLevel string) {
	config, err := ini.Load("../../.env")
	if err != nil {
		fmt.Printf("Error: configLoad: %v\n", err)
		return
	}
	logLevel = config.Section("logger").Key("LOG_LEVEL").MustString("")
	logPath = config.Section("logger").Key("LOG_PATH").MustString("")
	return
}

func TestLogger(t *testing.T) {

	logPath, logLevel := getConfig()

	logger := logger.NewLogger(logPath, logLevel)
	defer logger.Close()

	// DEBUG LEVEL
	logger.Debug("Please write messages here!", nil) // msgキーへ登録される

	// INFO LEVEL ,log & stdout
	logger.Info("Please write messages here!", "value", "key")    // "key":"value"形式と表示される
	logger.Warn("Please write messages here!", time.Now(), "now") // 時刻追加

	// WARN LEVEL ,log & stdout
	logger.Warn("Please write messages here!", nil)

	// ERROR LEVEL ,log & stdout
	// logger.Error("Please write messages here!", zap.Error(doEverything()))  // xerrorと組み合わせて使うと、スタックトレース内容も出力できる
	logger.Error("Please write messages here!", doEverything(), "Error")  // xerrorと組み合わせて使うと、スタックトレース内容も出力できる
	logger.Error("Please write messages here!", doEverything2(), "Error") // xerrorでwrap使用しても、トレースできる

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
	logger.Error("Please write messages here!", user, "userObj")

	// ユーザー設定した配列ログ内容
	users := &users{
		{
			ID:  1,
			Msg: "testです",
		},
		{
			ID:  2,
			Msg: "testです",
		},
	}

	// user arrayを挿入する、レシーバ（※user objectのようにクロージャも利用可）
	logger.Error("Please write messages here!", users, "userArray")
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

type users []user

// zap.Arrayの引数
func (us *users) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for _, u := range *us {
		enc.AppendObject(&u)
	}
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
