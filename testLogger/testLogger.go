package testLogger

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/google/logger"
	"golang.org/x/xerrors"
)

const logPath = "./example.log"

// MakeLog1 : errorがfmt.Errorfのケース
func MakeLog1() {
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()

	defer logger.Init("LoggerExample", true, false, ioutil.Discard).Close()

	logger.Info("I'm about to do something!")
	if err := doSomething(); err != nil {
		logger.Errorf("Error running doSomething: %v", err)
		logger.Warningf("Warning running doSomething: %v", err)
		logger.Fatalf("Fatal running doSomething: %v", err)
		logger.Info("not reached here!")
	}
}

func doSomething() error {
	return fmt.Errorf("fmtを使ったエラーです")
}

// MakeLog2 :
func MakeLog2() {
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()

	defer logger.Init("LoggerExample", true, false, lf).Close()

	logger.Info("I'm about to do something!")
	if err := doSomething2(); err != nil {
		logger.Errorf("Error running doSomething: %+v", err)
		logger.Warningf("Warning running doSomething: %+v", err)
		logger.Fatalf("Fatal running doSomething: %+v", err)
		logger.Info("not reached here!")
	}
}

func doSomething2() error {
	return xerrors.New("xerrorsを使ったエラーです")
}
