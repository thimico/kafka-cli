package utils

import (
	"github.com/thimico/kafka-cli/log"
	"go.uber.org/zap"
	"os"
)

const (
	DefaultErrorExitCode = 1
)

func fatal(msg string, code int) {
	os.Exit(code)
}

func checkErr(err error, handleErr func(string, int)) {
	if err == nil {
		return
	}
	log.Error("checkErr", zap.Error(err))
	handleErr(err.Error(), DefaultErrorExitCode)
}

func CheckErr(err error) {
	checkErr(err, fatal)
}
