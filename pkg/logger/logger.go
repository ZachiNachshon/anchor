package logger

import (
	"fmt"
	"os"
)

type AnchorLogger struct {
	//logImpl *logrus.Logger
	//logImpl *fmt.
}

var logger *AnchorLogger

func init() {
	//logrusInstance := logrus.New()
	logger = &AnchorLogger{
		//logImpl: logrusInstance,
	}
}

func Info(message string) {
	//logger.logImpl.Info(args)
	fmt.Println(message)
}

func Infof(format string, a ...interface{}) {
	//logger.logImpl.Info(args)
	fmt.Printf(format+"\n", a)
}

func Fatal(message string) {
	//logger.logImpl.Info(args)
	fmt.Println(message)
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	//logger.logImpl.Info(args)
	fmt.Printf(format, a)
	os.Exit(1)
}

func PrintHeadline(headline string) {
	format := fmt.Sprintf("----------------------- %s -----------------------", headline)
	Info(format)
}

func PrintCompletion() {
	format := "\n    Done.\n"
	Info(format)
}
