package main

// ZMQScopeTarget Portée du message
type ZMQScopeTarget int

const (
	// ZMQInternal Communications internes
	ZMQInternal ZMQScopeTarget = iota
	// ZMQDrone Communications externes
	ZMQDrone ZMQScopeTarget = iota
)

// ZMQIdentificationRequest Identification pour la classification CZMQ
type ZMQIdentificationRequest struct {
	Name     string         `json:"name"`
	ZMQPort  int            `json:"zmq_port"`
	Scope    ZMQScopeTarget `json:"scope"`
	Internal bool           `json:"internal"`
}

// ZMQInternalSystems Liste des systèmes internes connus par défaut
type ZMQInternalSystems string

const (
	// ZOSMService Service de planification des cartes
	ZOSMService ZMQInternalSystems = "locuste.services.osm"
)

// ZMQDefinedFunc Noms des fonctions échangées Router <=> Dealer
type ZMQDefinedFunc string

// ZMQMessage Message envoyé entre les Dealers ZMQ
type ZMQMessage struct {
	Function ZMQDefinedFunc `json:"function"`
	Params   []interface{}  `json:"params"`
}
