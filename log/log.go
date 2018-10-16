package log

import (
	"fmt"
	"log"
	"os"
)

var debug = false
var std = log.New(os.Stdout, "", log.LstdFlags)

func EnableDebug() {
	debug = true
}


func Info(l interface{}) {
	info := fmt.Sprintf("[%s] %v", "INFO", l)
	std.Output(2, info)
}

func Infof(l string, a ...interface{}) {
	tmp := fmt.Sprintf(l, a...)
	info := fmt.Sprintf("[%s] %s", "INFO", tmp)
	std.Output(2, info)
}

func Error(l interface{}) {
	err := fmt.Sprintf("[%s] %v", "Error", l)
	log.Println(err)
}

func Errorf(format string, a ...interface{}) {
	tmp := fmt.Sprintf(format, a...)
	err := fmt.Sprintf("[%s] %s", "Error", tmp)
	log.Println(err)
}

func Debug(l interface{})  {
	if debug {
		debug := fmt.Sprintf("[%s] %v", "DEBUG", l)
		log.Println(debug)
	}
}

func Debugf(format string, a ...interface{}) {
	if debug {
		tmp := fmt.Sprintf(format, a...)
		debug := fmt.Sprintf("[%s] %s", "Error", tmp)
		log.Println(debug)
	}
}