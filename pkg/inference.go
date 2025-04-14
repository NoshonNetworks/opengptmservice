package pkg

import (
	"log"
)

type Logger struct {
	Message string
}

func (l *Logger) Log(message Message) {

	message := "hiii"
	log.Println(message)
}
