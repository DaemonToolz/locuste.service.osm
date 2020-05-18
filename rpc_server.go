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
func (_ *RPCRegistry) RequestStatuses(_ *struct{}, reply *map[Component]bool) error {
	for key := range GlobalStatuses {
		(*reply)[key] = GlobalStatuses[key]
	}
	return nil
}

// OnCommandSuccess En cas de succès de la commande précédemment envoyée
func (_ *RPCRegistry) OnCommandSuccess(args *DroneIdentifier, _ *struct{}) error {
	go OnCommandSuccess(*args)
	return nil
}

// UpdateAutopilot En cas de mise à jour du pilote automatique
func (_ *RPCRegistry) UpdateAutopilot(input *SchedulerSummarizedData, _ *struct{}) error {
	myInput := *input
	log.Println("Demande de mise à jour du pilote automatique", myInput)
	go SendUpdateToScheduler(myInput)
	return nil
}

// UpdateTarget En cas de mise à jour de la destination
func (_ *RPCRegistry) UpdateTarget(input *FlightCoordinate, _ *struct{}) error {
	target := *input
	go func(coordinates FlightCoordinate) {
		UpdateSchedulerTarget(coordinates)
		SendEventToScheduler(coordinates.Name, OnTargetDefined)
	}(target)
	return nil
}

// GetTarget En cas de demande de la dernière destination
func (_ *RPCRegistry) GetTarget(args *DroneIdentifier, reply *FlightCoordinate) error {
	*reply = GetSchedulerTarget(args.Name)
	return nil
}

// GetBoundaries Récupère les bords
func (_ *RPCRegistry) GetBoundaries(_ *struct{}, reply *FlightBounds) error {
	if streetDataSet.Boundaries == nil {
		*reply = FlightBounds{}
	} else {
		*reply = *streetDataSet.Boundaries
	}
	return nil
}

// RestartModule Demande de redémarrage d'un module
func (_ *RPCRegistry) RestartModule(args *string, _ *struct{}) error {
	log.Println("Demande de redémarrage du module :", *args)
	return nil
}
