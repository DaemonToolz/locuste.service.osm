package main

import (
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("[ERROR] - %s: %s", msg, err)
	}
}
