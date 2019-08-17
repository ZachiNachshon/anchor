package logger

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"os"
)

type AnchorLogger struct {
	//logImpl *logrus.Logger
	//logImpl *fmt.
}

type HeadlineContext string

const DockerHeadline HeadlineContext = "Docker"
const ClusterHeadline HeadlineContext = "Cluster"

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
	//fmt.Println(message)
	fmt.Printf(Sprintf(Bold(Red("%v")), message))
	os.Exit(1)
}

func Fatalf(format string, a ...interface{}) {
	//logger.logImpl.Info(args)
	fmt.Printf(Sprintf(Bold(Red(format)), a))
	os.Exit(1)
}

func PrintHeadline(context HeadlineContext, action string) {
	var format = "\n[%s: %s]"
	switch context {
	case DockerHeadline:
		prefix := Sprintf(Bold(Blue("\n[%s: ")), context)
		suffix := Sprintf(Bold(Yellow("%s]")), action)
		format = prefix + suffix
		break
	case ClusterHeadline:
		prefix := Sprintf(Bold(Magenta("\n[%s: ")), context)
		suffix := Sprintf(Bold(Yellow("%s]")), action)
		format = prefix + suffix
		break
	}
	Info(format)
}

func PrintCompletion() {
	format := Sprintf(Bold(Green("\n    Done.\n")))
	Info(format)
}

func PrintCommandHeader(header string) {
	format := Sprintf(Bold(Cyan("\n ==> %v...\n")), header)
	Info(format)
}

func PrintWarning(msg string) {
	wrap := fmt.Sprintf(`
==========
IMPORTANT: %v
==========`, msg)
	format := Sprintf(Bold(Red("%v")), wrap)
	Info(format)
}
