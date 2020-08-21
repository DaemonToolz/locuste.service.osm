package main

import (
	"fmt"
	"sync"
)

var zmqServerMapMutex sync.Mutex
var zmqServerMap map[ZMQDefinedFunc]interface{}

func init() {
	zmqServerMap := make(map[ZMQDefinedFunc]interface{})

	zmqServerMap[ZFNRequestStatuses] = ZRequestStatuses
	zmqServerMap[ZFNNotifyScheduler] = ZNotifyScheduler
	zmqServerMap[ZFNUpdateAutopilot] = ZUpdateAutopilot
	zmqServerMap[ZFNOnHomeChanged] = ZOnHomeChanged
	zmqServerMap[ZFNFetchBoundaries] = ZFetchBoundaries
	zmqServerMap[ZFNUpdateTarget] = ZUpdateTarget
	zmqServerMap[ZFNUpdateFlyingStatus] = ZUpdateFlyingStatus
	zmqServerMap[ZFNSendGoHomeCommandTo] = ZSendGoHomeCommandTo
	zmqServerMap[ZFNSendTakeoffCommandTo] = ZSendTakeoffCommandTo
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
	// ZFNRequestStatusReply Fonction réponse de RequestStatuses
	ZFNRequestStatusReply ZMQDefinedFunc = "RequestStatusesReply"
)

// ZRegister Enregistrer le processus ZMQ
func ZRegister() *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNRegister,
		Params:   make([]interface{}, 0),
	}
}

// ZRDisconnect Désenregistre le process associé à une file ZMQ
func ZRDisconnect() *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNDisconnect,
		Params:   make([]interface{}, 0),
	}
}

// ZRequestStatusReply Message de réponse à la requête RequestStatus
func ZRequestStatusReply(param []interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNRequestStatusReply,
		Params:   param,
	}
}

// ZDefineBoundaries Définition des bordures de la carte
func ZDefineBoundaries(param []interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNDefineBoundaries,
		Params:   param,
	}
}

// ZSendCoordinates Envoi des coordonnées (informations automate)
func ZSendCoordinates(param []interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNSendCoordinates,
		Params:   param,
	}
}

// ZDefineTarget Défintion de la cible (à envoyer au drone)
func ZDefineTarget(param []interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNDefineTarget,
		Params:   param,
	}
}

// ZOnUpdateAutopilot Mise à jour de l'état de l'automate GO
func ZOnUpdateAutopilot(param []interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNOnUpdateAutopilot,
		Params:   param,
	}
}

// ZOnFlyingStatusUpdate Mise à jour des infos de vol
func ZOnFlyingStatusUpdate(param []interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNOnFlyingStatusUpdate,
		Params:   param,
	}
}

// ZServerShutdown Arrêt du serveur ZMQ
func ZServerShutdown(param []interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNServerShutdown,
		Params:   param,
	}
}

// ZSendCommand Envoi d'une commande
func ZSendCommand(param []interface{}) *ZMQMessage {
	return &ZMQMessage{
		Function: ZFNSendCommand,
		Params:   param,
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

)

// ZRequestStatuses Demande d'envoi des derniers status
func ZRequestStatuses() {
	outputState := make([]interface{}, 1)
	outputState[0] = GlobalStatuses
	// Request -> Reply
	go SendToZMQMessageChannel(ZRequestStatusReply(outputState))
	trace(callSuccess)
}

// ZNotifyScheduler Notification de la Finite State Machine (Scheduler)
func ZNotifyScheduler() {
	trace(callSuccess)
}

// ZUpdateAutopilot Réception des infos de MàJ du pilote
func ZUpdateAutopilot() {
	trace(callSuccess)
}

// ZOnHomeChanged Changement de la position "HOME"
func ZOnHomeChanged() {
	trace(callSuccess)
}

// ZFetchBoundaries Envoi des informations de géo-localisation
func ZFetchBoundaries() {
	trace(callSuccess)
}

// ZUpdateTarget Mise à jour de la cible de déplacement
func ZUpdateTarget() {
	trace(callSuccess)
}

// ZUpdateFlyingStatus Mise à jour des états de vol
func ZUpdateFlyingStatus() {
	trace(callSuccess)
}

// ZSendGoHomeCommandTo Processus de retour maison
func ZSendGoHomeCommandTo() {
	trace(callSuccess)
}

// ZSendTakeoffCommandTo Processus de décollage automatisé
func ZSendTakeoffCommandTo() {
	trace(callSuccess)
}

// #endregion Function Host
