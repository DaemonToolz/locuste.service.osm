package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// RPCRegistry Informations enregistrées dans le composant RPC
type RPCRegistry struct{}

var flightSchedulerRPC *RPCRegistry

func initRemoteProcedureCall() {
	listener, err := net.Listen("tcp", appConfig.rpcSchedulerPort())
	if err != nil {
		failOnError(err, "Couldn't initialize the RPC listener")
	}

	if flightSchedulerRPC == nil {
		flightSchedulerRPC = &RPCRegistry{}
		rpc.Register(flightSchedulerRPC)
		rpc.HandleHTTP()
	}
	log.Println("Ouverture des ports HTTP pour le processus RPC", listener.Addr().(*net.TCPAddr).Port)
	http.Serve(listener, nil)
}

// RestartRPCServer Redémarrer le serveur RPC
func RestartRPCServer() {
	initConfiguration()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				AddOrUpdateStatus(SchedulerRPCServer, false)
			}
		}()
		AddOrUpdateStatus(SchedulerRPCServer, true)
		initRemoteProcedureCall()
		AddOrUpdateStatus(SchedulerRPCServer, false)
		OnServerShutdown()
		log.Println("Arrêt du serveur RPC")
	}()

}

// RequestStatuses Fonction pour requêter les statuts des modules internes
func (*RPCRegistry) RequestStatuses(_ *struct{}, reply *map[Component]bool) error {
	for key := range GlobalStatuses {
		(*reply)[key] = GlobalStatuses[key]
	}
	return nil
}

// OnCommandSuccess En cas de succès de la commande précédemment envoyée
func (*RPCRegistry) OnCommandSuccess(args *CommandIdentifier, _ *struct{}) error {
	go OnCommandSuccess(*args)
	return nil
}

// UpdateAutopilot En cas de mise à jour du pilote automatique
func (*RPCRegistry) UpdateAutopilot(input *SchedulerSummarizedData, _ *struct{}) error {
	myInput := *input
	log.Println("Demande de mise à jour du pilote automatique", myInput)
	go SendUpdateToScheduler(myInput)
	return nil
}

// UpdateTarget En cas de mise à jour de la destination
func (*RPCRegistry) UpdateTarget(input *FlightCoordinate, _ *struct{}) error {
	target := *input
	go func(coordinates FlightCoordinate) {
		UpdateSchedulerTarget(coordinates)
		SendEventToScheduler(coordinates.Name, OnTargetDefined)
	}(target)
	return nil
}

// GetTarget En cas de demande de la dernière destination
func (*RPCRegistry) GetTarget(args *DroneIdentifier, reply *FlightCoordinate) error {
	*reply = GetSchedulerTarget(args.Name)
	return nil
}

// GetBoundaries Récupère les bords
func (*RPCRegistry) GetBoundaries(_ *struct{}, reply *FlightBounds) error {
	if streetDataSet.Boundaries == nil {
		*reply = FlightBounds{}
	} else {
		*reply = *streetDataSet.Boundaries
	}
	return nil
}

// FlyingStatusUpdate mise à jour de l'état du drone (en vol)
func (*RPCRegistry) FlyingStatusUpdate(data *DroneFlyingStatusMessage, _ *struct{}) error {
	UpdateFlyingStatus(data.Name, data.Status)

	return nil
}

// RestartModule Demande de redémarrage d'un module
// Stratégie aggressive comparé au pooling via le ping
func (*RPCRegistry) RestartModule(args *string, _ *struct{}) error {
	log.Println("Demande de redémarrage du module :", *args)
	// TODO : Logique à implémenter
	return nil
}

// SendGoHomeCommandTo Demander une commande "atterrissage" au drone nommé
func (*RPCRegistry) SendGoHomeCommandTo(name *string, _ *struct{}) error {
	SendEventToScheduler(*name, OnGoHome) // Implémenter la logique
	return nil
}

// SendTakeoffCommandTo Demander une commande "décollage" au drone nommé
func (*RPCRegistry) SendTakeoffCommandTo(name *string, _ *struct{}) error {
	SendEventToScheduler(*name, OnTakeOff) // Implémenter la logique
	return nil
}
