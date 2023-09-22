package util

import (
	"fmt"
	"log"
)

var Debug bool

func init() {

}

// LogDebug logs a message to stdout if debugging is enabled
func LogDebug(msg string, args ...interface{}) {
	//if Debug {
	var message string
	if len(args) == 0 {
		message = msg
	} else {
		message = fmt.Sprintf(msg, args...)
	}
	log.Println(message)
	//}
}
