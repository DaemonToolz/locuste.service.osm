package main

import (
	"fmt"
	"log"
)

// ExtractDroneNames Récupère les noms des drones
func ExtractDroneNames() []string {
	droneNames := make([]string, 0)

	for _, drone := range drones.Drones {
		droneNames = append(droneNames, fmt.Sprintf("ANAFI_%s", drone.IpAddress))
	}

	return droneNames
}

// ModuleRestart Redémarrage d'un module
func ModuleRestart(module Module) {
	log.Println("Module à redémarrer : ", module)
	CallModuleRestart(Component(module.System + "." + module.SubSystem))
}
