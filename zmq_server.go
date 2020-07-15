package main

import (
	"github.com/zeromq/goczmq"
)

// #region Section CZMQ
var zmqBrainSocket *goczmq.Sock
var zmqCmdChannel chan interface{}

// #endregion Section CZMQ

func initZMQ() {
	zmqCmdChannel = make(chan interface{})
}
