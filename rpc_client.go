package main

import (
	"log"
	"net/rpc"
	"os"
	"reflect"
	"time"
)

// Args Argument d'identification (renommer)
type Args struct {
	PId       int
	Component Component
}

var client *rpc.Client
var myself Args

var pulse *time.Ticker
var stopCondition chan bool
var lastStatuses map[Component]bool

/* NullArgType Structure leurre pour pour envoyer une struct "nil"
Possiblement remplacer les connexions RPC par du gRPC
*/
type NullArgType struct{}

// RPCNullArg Objet null commun à tous les appels RPC
var RPCNullArg NullArgType

// ModuleToRestart Module à redémarrer (envoyé par l'unité de contrôle, ou "brain")
var ModuleToRestart string

func initRPCClient() {
	ModuleToRestart = ""
	RPCNullArg = NullArgType{}
	myself = Args{
		PId:       os.Getpid(), // Informations à transmettre: notre process + notre module
		Component: SchedulerRPC,
	}

	pulse = time.NewTicker(1 * time.Second)
	stopCondition = make(chan bool)
	openConnection()
	go ping()
	log.Println("Connectiques RPC initialisés")
}

func ping() {
	for {
		select {
		case <-stopCondition:
			log.Println("Connectiques RPC arrêtées")

			return
		case <-pulse.C:
			if client != nil {
				accessCall := client.Go("RPCRegistry.Ping", &RPCNullArg, &ModuleToRestart, nil)
				replyCall := <-accessCall.Done

				if client == nil {
					log.Println("La connexion n'était pas initialisée")
					openConnection()
				} else if replyCall.Error == rpc.ErrShutdown || reflect.TypeOf(replyCall.Error) == reflect.TypeOf((*rpc.ServerError)(nil)).Elem() {
					log.Println("Une erreur liée au serveur a été remonté")
					log.Println(replyCall.Error)
					openConnection()
				}

			} else {
				openConnection()
			}

			if ModuleToRestart != "" {
				CallModuleRestart(Component(ModuleToRestart))
			}

			ModuleToRestart = ""
		}
	}
}

func openConnection() *rpc.Client {
	initConfiguration()
	var err error
	client, err = rpc.DialHTTP("tcp", appConfig.rpcListenUri())
	if err != nil {
		AddOrUpdateStatus(SchedulerRPC, false)
		failOnError(err, "couldn't connect to remote RPC server")
	} else {

		AddOrUpdateStatus(SchedulerRPC, true)
	}

	if client != nil {
		client.Go("RPCRegistry.Register", &myself, &RPCNullArg, nil)
	}
	return client
}

// Unregister On a effectué le Register auprès du service central, il faut appeler Unregister / Disconnect
func Unregister() {
	defer func() {
		if client != nil {
			defer client.Close()
		}
	}()

	if client != nil {
		client.Go("RPCRegistry.Disconnect", &myself, &RPCNullArg, nil)
	}
	AddOrUpdateStatus(SchedulerRPC, false)
}

// TransmitCoordinates Envoyer des nouvelles coordonnées à l'unité contrôle
func TransmitCoordinates(coordinate *DroneFlightCoordinates) {
	if client != nil {
		client.Go("RPCRegistry.SendCoordinates", coordinate, &RPCNullArg, nil)
	}
}

// TransmitTarget Transmet la cible recalculée
func TransmitTarget(coordinate *FlightCoordinate) {
	if client != nil {
		client.Go("RPCRegistry.DefineTarget", coordinate, &RPCNullArg, nil)
	}
}

// TransmitEdge Deprecated: Permet d'envoyer les liaisons à l'unité de contrôle
func TransmitEdge(coordinate *FlightCoordinate) {
	if client != nil {
		client.Go("RPCRegistry.DefineEdge", coordinate, &RPCNullArg, nil)
	}
}

// TransmitAutopilotUpdate Transmet les mises à jour du pilote automatique
func TransmitAutopilotUpdate(input *SchedulerSummarizedData) {
	if client != nil {
		client.Go("RPCRegistry.UpdateAutopilot", input, &RPCNullArg, nil)
	}
}

// TransmitEvent Transmet une commande
func TransmitEvent(command *DroneCommandMessage) {
	if client != nil {
		client.Go("RPCRegistry.RPCSendCommand", command, &RPCNullArg, nil)
	}
}

// TransmitBounds Transmet les bords de la carte
func TransmitBounds(boundaries *FlightBounds) {
	if client != nil {
		client.Go("RPCRegistry.DefineBoundaries", boundaries, &RPCNullArg, nil)
	}
}

// OnServerShutdown Appeler la procédure d'arrêt
func OnServerShutdown() {
	if client != nil {
		AddOrUpdateStatus(SchedulerRPCServer, false)
		client.Go("RPCRegistry.ServerShutdown", &RPCNullArg, &RPCNullArg, nil)
	}
}
