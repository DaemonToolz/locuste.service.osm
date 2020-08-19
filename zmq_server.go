package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/zeromq/goczmq"
)

// #region Section CZMQ
var zmqBrainSocket *goczmq.Sock
var zmqCmdChannel chan interface{}
var zmqBrainAccessMutex sync.Mutex
var zmqBrainCmdMutex sync.Mutex

// #endregion Section CZMQ

func initZMQ() {
	zmqCmdChannel = make(chan interface{})

	CreateZMQDealer(ZMQIdentificationRequest{
		Name:    string(ZOSMService),
		Scope:   ZMQInternal,
		ZMQPort: appConfig.OSMZmqPort,
	}, true)
}

func createTgtZMQDealer(who *goczmq.Sock, how *chan interface{}, modMutex *sync.Mutex, cmdMutex *sync.Mutex, request *ZMQIdentificationRequest) {
	var err error
	modMutex.Lock()
	defer modMutex.Unlock()

	cmdMutex.Lock()
	close((*how))
	cmdMutex.Unlock()

	who, err = goczmq.NewDealer(fmt.Sprintf("tcp://127.0.0.1:%d", request.ZMQPort)) // 5555 is Default ZMQ port

	if err != nil {
		failOnError(err, "CreateZMQDealer")
		AddOrUpdateStatus(Component(request.Name), false)

		cmdMutex.Lock()
		(*how) = make(chan interface{})
		cmdMutex.Unlock()
		AddOrUpdateStatus(Component(request.Name), false)

		go func(toWhom *goczmq.Sock, commChan *chan interface{}, lock *sync.Mutex, name string) {
			messageListenerLoop(toWhom, commChan, lock, name)
		}(who, how, modMutex, request.Name)
	}

	log.Println("Dealer ZeroMQ initialisé")
}

func messageListenerLoop(toWhom *goczmq.Sock, commChan *chan interface{}, lock *sync.Mutex, name string) {
	for data := range *commChan {
		SendZMQMessage(toWhom, lock, name, data)
	}
}

func messageReceiverLoop(toWhom *goczmq.Sock, commChan *chan interface{}, lock *sync.Mutex, name string) {
	for {
		msg, err := toWhom.RecvMessage()
		if err != nil {
			failOnError(err, "Error in messageReceiverLoop")
			DestroyZMQDealer(toWhom, lock, name) // On détruit tout, car un bug s'est présenté
		}
		var payload ZMQMessage
		err = json.Unmarshal(msg[0], payload)

		if err != nil {
			failOnError(err, "Error in messageReceiverLoop")
			DestroyZMQDealer(toWhom, lock, name) // On détruit tout, car un bug s'est présenté
		}

		callMappedZMQFunc(&payload)
	}
}

// SendZMQMessage Envoyer un message via ZMQ (not thread-safe)
func SendZMQMessage(toWhom *goczmq.Sock, lock *sync.Mutex, name string, payload interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)

			AddOrUpdateStatus(Component(name), false)

			DestroyZMQDealer(toWhom, lock, name) // On détruit tout, car un bug s'est présenté
		}
	}()

	jPayload, err := json.Marshal(&payload)
	if err != nil {
		failOnError(err, fmt.Sprintf("SendZMQMessage:%s", name))
		return
	}

	lock.Lock()
	_, erro := toWhom.Write([]byte(jPayload))
	lock.Unlock()
	if erro != nil {
		failOnError(erro, fmt.Sprintf("SendZMQMessage:%s", name))
	}

}

// DestroyZMQRouters Destruction de toutes les connectiques CZMQ
func DestroyZMQRouters() {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()
	DestroyZMQDealer(zmqBrainSocket, &zmqBrainAccessMutex, string(SchedulerRPCServer))

}

// DestroyZMQDealer Destruction d'un dealer CZMQ
func DestroyZMQDealer(who *goczmq.Sock, modMutex *sync.Mutex, name string) {

	AddOrUpdateStatus(Component(name), false)

	modMutex.Lock()
	defer modMutex.Unlock()
	if who != nil {
		who.Destroy()
	}

}

// CreateZMQDealer Création d'un Dealer CZMQ (rattaché au zmqSocket)
func CreateZMQDealer(request ZMQIdentificationRequest, internal bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	createTgtZMQDealer(zmqBrainSocket, &zmqCmdChannel, &zmqBrainAccessMutex, &zmqBrainCmdMutex, &request)
}
