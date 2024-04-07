package logger

import (
	"log"

	"github.com/fatih/color"
)

var red = color.New(color.FgRed)
var green = color.New(color.FgGreen)
var yellow = color.New(color.FgYellow)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func print(message ...interface{}) {
	log.Println(message...)
}

func Info(message ...interface{}) {
	print("[INFO]: ", yellow.Sprint(message...))
}

func Error(message ...interface{}) {
	print("[ERROR]: ", red.Sprint(message...))
}
func Success(message ...interface{}) {
	print("[SUCCESS]: ", green.Sprint(message...))
}
