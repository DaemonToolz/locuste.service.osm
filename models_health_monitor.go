package main

import "sync"

/*
	Regroupement des éléments purements dédiés à la partie GO
*/

var moduleMutex sync.Mutex

/* Component Composant logique de la brique logicielle
e.g. Connectiques RPC, machines à états
*/
type Component string

const (
	// SchedulerRPCServer Serveur RPC de l'ordonnanceur
	SchedulerRPCServer Component = "Scheduler.RPCServer"
	// SchedulerRPC Connexion RPC de l'ordonnanceur vers le serveur RPC de l'unité de contrôle
	SchedulerRPC Component = "Scheduler.BrainConnection"
	// SchedulerMapHandler Gestionnaire de cartes
	SchedulerMapHandler Component = "Scheduler.MapHandler"
	// SchedulerFlightManager Gestionnaire de vol / Pilotes automatiques / Ordonnanceurs de vols
	SchedulerFlightManager Component = "Scheduler.FlightManager"
)

// Module System - Sous-sytème
type Module struct {
	System    string `json:"system"`
	SubSystem string `json:"subsystem"`
}

// ModuleRestartMapper Mappeur global
var ModuleRestartMapper map[Component]interface{}

func initModuleRestartMapper() {
	ModuleRestartMapper = make(map[Component]interface{})
	AddOrUpdateModuleMapper(SchedulerRPCServer, RestartRPCServer)
	AddOrUpdateModuleMapper(SchedulerFlightManager, RestartSchedulers)
}

// AddOrUpdateModuleMapper Ajout d'une fonction de mise à jour
func AddOrUpdateModuleMapper(comp Component, function interface{}) {
	moduleMutex.Lock()
	ModuleRestartMapper[comp] = function
	moduleMutex.Unlock()
}

// CallModuleRestart Fonction pour redémarrer un module
func CallModuleRestart(comp Component) bool {
	moduleMutex.Lock()
	defer moduleMutex.Unlock()
	if _, ok := ModuleRestartMapper[comp]; ok {
		if result, ok := GlobalStatuses[comp]; !ok || (ok && !result) {
			ModuleRestartMapper[comp].(func())()
			return true
		}

	}
	return false
}
