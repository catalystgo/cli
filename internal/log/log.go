package log

import (
	"log"

	"github.com/fatih/color"
)

var (
	verbose bool

	prefixDegbug = color.CyanString("[DBG] ")
	prefixInfo   = color.GreenString("[INF] ")
	prefixWarn   = color.YellowString("[WAR] ")
	prefixError  = color.RedString("[ERR] ")
)

func SetVerbose(b bool) {
	verbose = b
}

func Debug(msg string) {
	if verbose {
		log.Println(prefixDegbug + msg)
	}
}

func Debugf(msg string, args ...interface{}) {
	if verbose {
		log.Printf(prefixDegbug+msg, args...)
	}
}

func Info(msg string) {
	log.Println(prefixInfo + msg)
}

func Infof(msg string, args ...interface{}) {
	log.Printf(prefixInfo+msg, args...)
}

func Warn(msg string) {
	log.Println(prefixWarn + msg)
}

func Warnf(msg string, args ...interface{}) {
	log.Printf(prefixWarn+msg, args...)
}

func Error(msg string) {
	log.Println(prefixError + msg)
}

func Errorf(msg string, args ...interface{}) {
	log.Printf(prefixError+msg, args...)
}
