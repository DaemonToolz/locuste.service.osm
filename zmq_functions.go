package main

import (
	"sync"
)

// ZMQDefinedFunc Noms des fonctions échangées Router <=> Dealer
type ZMQDefinedFunc string

// ZMQMessage Message envoyé entre les Dealers ZMQ
type ZMQMessage struct {
	Function ZMQDefinedFunc `json:"function"`
	Params   []interface{}  `json:"params"`
}

// ZMQComponents Composants enregistrés (Avec ProcessID)
var ZMQComponents map[Component]int
var zmqProcessMutex sync.Mutex

func init() {
	ZMQComponents = make(map[Component]int)
}

// #region Function Client
const (
	// ZFNRegister Anciennement Register (RPC)
	ZFNRegister ZMQDefinedFunc = "Register"
	// ZFNDisconnect Anciennement Disconnect (RPC)
	ZFNDisconnect ZMQDefinedFunc = "Disconnect"
	// ZFNDisconnect Anciennement DefineBoundaries (RPC)
	ZFNDefineBoundaries ZMQDefinedFunc = "DefineBoundaries"
	// ZFNSendCoordinates Anciennement SendCoordinates (RPC)
	ZFNSendCoordinates ZMQDefinedFunc = "SendCoordinates"
	// ZFNDefineTarget Anciennement DefineTarget (RPC)
	ZFNDefineTarget ZMQDefinedFunc = "DefineTarget"
	// ZFNOnUpdateAutopilot Anciennement UpdateAutopilot (RPC)
	ZFNOnUpdateAutopilot ZMQDefinedFunc = "OnUpdateAutopilot"
	// ZFNOnFlyingStatusUpdate Anciennement OnFlyingStatusUpdate (RPC)
	ZFNOnFlyingStatusUpdate ZMQDefinedFunc = "OnFlyingStatusUpdate"
	// ZFNServerShutdown Anciennement ServerShutdown (RPC)
	ZFNServerShutdown ZMQDefinedFunc = "ServerShutdown"
	// ZFNSendCommand Anciennement SendCommand (RPC)
	ZFNSendCommand ZMQDefinedFunc = "SendCommand"
)

func addOrUpdateZMQProcess(cpt Component, pid int) {
	zmqProcessMutex.Lock()
	defer zmqProcessMutex.Unlock()
	ZMQComponents[cpt] = pid
}

func deleteZMQProcess(cpt Component) {
	zmqProcessMutex.Lock()
	defer zmqProcessMutex.Unlock()
	delete(ZMQComponents, cpt)
}

// #endregion Function Client

// #region Function Host
const (
	// ZFNRequestStatuses Anciennement RequestStatuses (RPC)
	ZFNRequestStatuses ZMQDefinedFunc = "RequestStatuses"
	// ZFNNotifyScheduler Anciennement OnCommandSuccess (RPC)
	ZFNNotifyScheduler ZMQDefinedFunc = "OnCommandSuccess"
	// ZFNUpdateAutopilot Anciennement UpdateAutopilot (RPC)
	ZFNUpdateAutopilot ZMQDefinedFunc = "UpdateAutopilot"
	// ZFNOnHomeChanged Anciennement OnHomeChanged (RPC)
	ZFNOnHomeChanged ZMQDefinedFunc = "OnHomeChanged"
	// ZFNFetchBoundaries Anciennement GetBoundaries (RPC)
	ZFNFetchBoundaries ZMQDefinedFunc = "GetBoundaries"
	// ZFNUpdateTarget Anciennement UpdateTarget (RPC)
	ZFNUpdateTarget ZMQDefinedFunc = "UpdateTarget"
	// ZFNUpdateFlyingStatus Anciennement FlyingStatusUpdate (RPC)
	ZFNUpdateFlyingStatus ZMQDefinedFunc = "FlyingStatusUpdate"
	// ZFNSendGoHomeCommandTo Anciennement SendGoHomeCommandTo (RPC)
	ZFNSendGoHomeCommandTo ZMQDefinedFunc = "SendGoHomeCommandTo"
	// ZFNSendTakeoffCommandTo Anciennement SendTakeoffCommandTo (RPC)
	ZFNSendTakeoffCommandTo ZMQDefinedFunc = "SendTakeoffCommandTo"

	// Reply sections - From client

	// ZFNRequestStatusReply Fonction réponse de RequestStatuses
	ZFNRequestStatusReply ZMQDefinedFunc = "RequestStatusesReply"
)

// #endregion Function Host
