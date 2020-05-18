package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
)

func loadStateMachine() StateMachineJSON {
	var localFsm StateMachineJSON
	configFile, err := os.Open("./config/state_machine.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&localFsm)
	return localFsm
}

// InitStateMachine Initialisation de la machine à état
func (sd *SchedulerData) InitStateMachine() {
	sd.FSM.Init(sd.OnUpdateChannel)
	for _, jsonState := range loadStateMachine().StateMachine {
		state := State{
			Name:        jsonState.Name,
			Event:       Event(jsonState.Event),
			Description: jsonState.Description,
		}

		state.Init()

		for _, outcome := range jsonState.Outcomes {
			state.SetOutcome(FuncResult(outcome.Result), Event(outcome.NextState))
		}

		for _, income := range jsonState.Incomes {
			state.SetIncome(Event(income))
		}

		var pre interface{} = nil
		var call interface{} = nil
		var post interface{} = nil

		if jsonState.Callbacks.Precall != "" {
			pre = reflect.ValueOf(sd).MethodByName(jsonState.Callbacks.Precall).Interface()
		}
		call = reflect.ValueOf(sd).MethodByName(jsonState.Callbacks.Call).Interface()

		if jsonState.Callbacks.Postcall != "" {
			post = reflect.ValueOf(sd).MethodByName(jsonState.Callbacks.Postcall).Interface()
		}

		state.SetCallBacks(pre, call, post)
		sd.FSM.AddState(state)
	}

}

// DefineStartingPoint Définition du point de départ
func (sd *SchedulerData) DefineStartingPoint() FuncResult {
	log.Println("Le point de décollage / origine a changé: ", sd.DroneName)
	sd.OperationIndex = 0
	sd.ReadInitialDroneStatus()
	if sd.Statuses.IsReady {
		return OnSuccess
	} else {
		return OnError
	}
}

// SetNewRoute Définition d'une nouvelle route
func (sd *SchedulerData) SetNewRoute() FuncResult {
	if !sd.Statuses.IsManual {
		log.Println("Une nouvelle route est disponible : ", sd.DroneName)
		if sd.PrepareNextCoordinate() {
			TransmitCoordinates(sd.CurrentInstruction)
			return UpdateRequired
		}
	}

	return OnSuccess
}

// SetPositionReached Position atteinte
func (sd *SchedulerData) SetPositionReached() FuncResult {
	if !sd.Statuses.IsManual {
		log.Println("Position atteinte : ", sd.DroneName)
		if sd.OperationIndex < len(*sd.Route) {
			sd.OperationIndex++
			if sd.PrepareNextCoordinate() {
				TransmitCoordinates(sd.CurrentInstruction)
				return UpdateRequired
			} else {
				return OnTerminated
			}
		}
	}
	return OnSuccess
}

// SetManual Passage en mode manuel
func (sd *SchedulerData) SetManual() FuncResult {
	log.Println("Passage en mode manuel : ", sd.DroneName)
	sd.Statuses.IsManual = true
	return OnSuccess
}

// SetAutomatic Passage en mode automatique
func (sd *SchedulerData) SetAutomatic() FuncResult {
	log.Println("Passage en mode automatique : ", sd.DroneName)
	sd.Statuses.IsManual = false
	return OnSuccess
}

// SetSimulation Passage en mode simulation
func (sd *SchedulerData) SetSimulation() FuncResult {
	log.Println("Passage en mode Simulation : ", sd.DroneName)
	sd.Statuses.IsSimulated = true
	sd.OperationIndex = 0
	return OnSuccess
}

// SetAutopilotOn Activation du pilote automatique
func (sd *SchedulerData) SetAutopilotOn() FuncResult {
	log.Println("Activation de l'autopilote : ", sd.DroneName)
	sd.Statuses.IsActive = true
	return OnSuccess
}

// SetAutopilotOff Désactivation du pilote automatique
func (sd *SchedulerData) SetAutopilotOff() FuncResult {
	log.Println("Désactivation de l'autopilote : ", sd.DroneName)
	sd.Statuses.IsActive = false
	return OnSuccess
}

// SetNormal Passage en mode normal
func (sd *SchedulerData) SetNormal() FuncResult {
	log.Println("Passage en mode standard : ", sd.DroneName)
	sd.Statuses.IsSimulated = false
	sd.OperationIndex = 0
	return OnSuccess
}

// SetDestinationReached Indique que la destination est atteinte
func (sd *SchedulerData) SetDestinationReached() FuncResult {
	log.Println("Navigation terminée :", sd.DroneName)
	return OnSuccess
}

// SetTargetDefined Définition d'une nouvelle cible
func (sd *SchedulerData) SetTargetDefined() FuncResult {
	log.Println("Navigation démarrée :", sd.DroneName)
	origin := GetSchedulerTarget(sd.DroneName)
	log.Println("Ancienne cible de navigation récupérée :", sd.DroneName, origin)

	target := DefineClosestStartingPoint(origin)
	sd.Target = target
	log.Println("Nouvelle cible de navigation définie :", sd.DroneName, target)

	toTransmit := &FlightCoordinate{
		Name: sd.DroneName,
		Lat:  target.Lat,
		Lon:  target.Lng,
	}

	UpdateSchedulerTarget(*toTransmit)
	TransmitTarget(toTransmit)
	return OnSuccess
}

// SetUpdate Appelle de la mise à jour des états via la machine à état
func (sd *SchedulerData) SetUpdate() FuncResult {
	log.Println("Autopilote pour ", sd.DroneName, " mis à jour")
	UpdateMapStatus(sd.DroneName, *sd)
	return OnSuccess
}

// UpdateMapStatus Appelle de la mise à jour des états via un appel direct
func (sd *SchedulerData) UpdateMapStatus() {
	sd.SetUpdate()
}
