package util

import (
	"fmt"
	"log"
)

// HandleError logs an error message to stderr if err is not nil
func HandleError(err interface{}, msg string, param ...any) {
	message := fmt.Sprintf(msg, param...)
	if err != nil {
		log.Println("Error: ", message, err)
	}
}

// HandleFatalError logs an error message to stderr and exits if err is not nil
func HandleFatalError(err interface{}, msg string, param ...any) {
	message := fmt.Sprintf(msg, param...)
	if err != nil {
		log.Fatalln("Error: ", message, err)
	}
}
