package main

import (
	"log"
	"runtime"
)

// #region Pre-recorded messages
const (
	callFailure string = "Echec de l'appel"
	callSuccess string = "Réussite de l'appel"
)

// #endregion Pre-recorded messages

func failOnError(err error, msg string) {
	if err != nil {
		log.Printf("[ERROR] - %s: %s", msg, err)
	}
}

func trace(desc string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	log.Printf("%s - %s\n", frame.Function, desc)
}
