package utils

import (
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func Info(message string) {
	// Callers(1, 1)로 호출된 파일명, 라인 번호, 함수명을 가져옴
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		logrus.Error("Failed to get caller info")
		return
	}

	// 함수명을 가져옴
	funcName := strings.Split(runtime.FuncForPC(pc).Name(), "/")[3]

	// 파일명, 라인 번호, 함수명을 포함한 로그 메시지 출력
	logrus.WithFields(logrus.Fields{
		"file":     file,
		"line":     line,
		"function": funcName,
	}).Info(message)
}

func Error(message string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		logrus.Error("Failed to get caller info")
		return
	}

	funcName := runtime.FuncForPC(pc).Name()

	logrus.WithFields(logrus.Fields{
		"file":     file,
		"line":     line,
		"function": funcName,
	}).Error(message)
}

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
}