package main

import (
	"github.com/kit/cmd"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	logrus.SetFormatter(&logrus.TextFormatter{})

	//Only logrus the warning severity or above.
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	cmd.NewCmdRoot().Execute()
}
