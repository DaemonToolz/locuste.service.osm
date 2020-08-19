package main

import (
	"fmt"
	"sync"
)

var zmqServerMapMutex sync.Mutex
var zmqServerMap map[ZMQDefinedFunc]interface{}

func init() {
	zmqServerMap := make(map[ZMQDefinedFunc]interface{})

	
	zmqServerMap[ZFNRequestStatuses] = nil
	zmqServerMap[ZFNNotifyScheduler] = nil
	zmqServerMap[ZFNUpdateAutopilot] = nil
	zmqServerMap[ZFNOnHomeChanged] = nil
	zmqServerMap[ZFNFetchBoundaries] = nil
	zmqServerMap[ZFNUpdateTarget] = nil
	zmqServerMap[ZFNUpdateFlyingStatus] = nil
	zmqServerMap[ZFNSendGoHomeCommandTo] = nil
	zmqServerMap[ZFNSendTakeoffCommandTo] = nil
	zmqServerMap[ZFNRequestStatusReply] = nil
	
}

func callMappedZMQFunc(msg *ZMQMessage) {
	if msg != nil {
		if _, ok := zmqServerMap[msg.Function]; !ok {
			trace(fmt.Sprintf("%s : %s", callFailure, "Méthode inconnue"))
		}
		zmqServerMapMutex.Lock()
		defer func() {
			zmqServerMapMutex.Unlock()
			if r := recover(); r != nil {
				trace(fmt.Sprintf("%s : %s", callFailure, r))
			}
		}()

		zmqServerMap[msg.Function].(func(*[]interface{}))(&msg.Params)
	} else {
		trace(fmt.Sprintf("%s : %s", callFailure, "Message reçu malformé"))
	}
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

// ZRegister Enregistrer le processus ZMQ
func ZRegister(params *[]interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNRegister,
		Params:   make([]interface{}, 0),
	}
}

// ZRDisconnect Désenregistre le process associé à une file ZMQ
func ZRDisconnect(params *[]interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNDisconnect,
		Params:   make([]interface{}, 0),
	}
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
