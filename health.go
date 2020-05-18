package main;

import "sync"

var statusMutex sync.Mutex
// GlobalStatuses Récupère l'état de fonctionnement des composants
var GlobalStatuses map[Component]bool


func initHealthMonitor(){
	GlobalStatuses = make(map[Component]bool)
}

// AddOrUpdateStatus Met à jour l'information d'un composant
func AddOrUpdateStatus(component Component, isOnline bool){
	statusMutex.Lock()
	GlobalStatuses[component] = isOnline;
	statusMutex.Unlock();
}

