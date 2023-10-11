package util

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

var Verbose bool
var Debug bool

func init() {

}

// Prints a message to stdout
func Log(msg string, args ...interface{}) {
	var message string

	Verbose = viper.GetBool("verbose")

	if Verbose {
		if len(args) == 0 {
			message = msg
		} else {
			message = fmt.Sprintf(msg, args...)
		}
		log.Println(message)
	}
}

// LogDebug logs a message to stdout if debugging is enabled
func LogDebug(msg string, args ...interface{}) {

	Debug = viper.GetBool("debug")

	if Debug {
		var message string
		if len(args) == 0 {
			message = msg
		} else {
			message = fmt.Sprintf(msg, args...)
		}

		log.Println(message)
	}
}
