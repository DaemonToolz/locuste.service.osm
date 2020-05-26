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
	sd.IntermediateCommand = NoCommand
	sd.LastCommand = NoCommand
	if sd.Statuses.IsReady {
		sd.LastCommand = GoTo
		if !sd.Statuses.IsManual {
			if sd.DroneFlyingStatus.IsLanded {
				sd.IntermediateCommand = TakeOff
			}
			if !sd.DroneFlyingStatus.IsGoingHome && sd.DroneFlyingStatus.IsMoving {
				sd.IntermediateCommand = Stop
			}

			if sd.IntermediateCommand != NoCommand {
				return MoveInterruptRequired
			}
		}

		return OnSuccess
	}
	return OnError
}

// SendIntermediateCommand On envoi la commande intermédiaire avant la commande finale
func (sd *SchedulerData) SendIntermediateCommand() FuncResult {
	payload := &DroneCommandMessage{
		Name:   sd.IntermediateCommand,
		Target: sd.DroneName,
	}
	TransmitEvent(payload)
	return OnSuccess
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

	sd.LastCommand = NoCommand
	return OnSuccess
}

// SetPositionReached Position atteinte
func (sd *SchedulerData) SetPositionReached() FuncResult {
	if !sd.Statuses.IsManual {
		log.Println("Position atteinte : ", sd.DroneName)
		if sd.Route != nil && sd.OperationIndex < len(*sd.Route) {
			sd.OperationIndex++
			if sd.PrepareNextCoordinate() && sd.IntermediateCommand == NoCommand {
				TransmitCoordinates(sd.CurrentInstruction)
				return UpdateRequired
			}
			return OnTerminated

		}
	}
	return OnSuccess
}

// SetManual Passage en mode manuel
func (sd *SchedulerData) SetManual() FuncResult {
	log.Println("Passage en mode manuel : ", sd.DroneName)
	sd.Statuses.IsManual = true
	sd.LastCommand = NoCommand
	sd.IntermediateCommand = NoCommand
	return OnSuccess
}

// SetAutomatic Passage en mode automatique
func (sd *SchedulerData) SetAutomatic() FuncResult {
	log.Println("Passage en mode automatique : ", sd.DroneName)
	sd.Statuses.IsManual = false
	sd.LastCommand = NoCommand
	sd.IntermediateCommand = NoCommand
	return OnSuccess
}

// SetSimulation Passage en mode simulation
func (sd *SchedulerData) SetSimulation() FuncResult {
	log.Println("Passage en mode Simulation : ", sd.DroneName)
	sd.Statuses.IsSimulated = true
	sd.OperationIndex = 0
	sd.LastCommand = NoCommand
	sd.IntermediateCommand = NoCommand
	return OnSuccess
}

// SetAutopilotOn Activation du pilote automatique
func (sd *SchedulerData) SetAutopilotOn() FuncResult {
	log.Println("Activation de l'autopilote : ", sd.DroneName)
	sd.Statuses.IsActive = true
	sd.LastCommand = NoCommand
	sd.IntermediateCommand = NoCommand
	return OnSuccess
}

// SetAutopilotOff Désactivation du pilote automatique
func (sd *SchedulerData) SetAutopilotOff() FuncResult {
	log.Println("Désactivation de l'autopilote : ", sd.DroneName)
	sd.Statuses.IsActive = false
	sd.LastCommand = NoCommand
	sd.IntermediateCommand = NoCommand
	return OnSuccess
}

// SetNormal Passage en mode normal
func (sd *SchedulerData) SetNormal() FuncResult {
	log.Println("Passage en mode standard : ", sd.DroneName)
	sd.Statuses.IsSimulated = false
	sd.OperationIndex = 0
	sd.LastCommand = NoCommand
	sd.IntermediateCommand = NoCommand
	return OnSuccess
}

// SetDestinationReached Indique que la destination est atteinte
func (sd *SchedulerData) SetDestinationReached() FuncResult {
	log.Println("Navigation terminée :", sd.DroneName)
	sd.LastCommand = NoCommand
	sd.IntermediateCommand = NoCommand
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

// OnFlyingStateReceived Appelle de la mise à jour des états via la machine à état
// A voir pour mettre en place une fonction Pre-call
func (sd *SchedulerData) OnFlyingStateReceived() FuncResult {
	log.Println("Informations de vol pour ", sd.DroneName, " mis à jour")
	switch GetFlyingStatus(sd.DroneName) {
	case NoStatus:
		sd.DroneFlyingStatus.IsLanded = false
		sd.DroneFlyingStatus.IsHovering = false
		sd.DroneFlyingStatus.IsMoving = false
		sd.DroneFlyingStatus.IsPreparing = false
		sd.DroneFlyingStatus.IsGoingHome = false

	case Landed:
		sd.DroneFlyingStatus.IsLanded = true
		sd.DroneFlyingStatus.IsHovering = false
		sd.DroneFlyingStatus.IsMoving = false
		sd.DroneFlyingStatus.IsPreparing = false
		sd.DroneFlyingStatus.IsGoingHome = false

	case Flying:
		sd.DroneFlyingStatus.IsLanded = false
		sd.DroneFlyingStatus.IsHovering = false
		sd.DroneFlyingStatus.IsMoving = true
		sd.DroneFlyingStatus.IsPreparing = false
		// Peut être l'instruction Go Home, dans ce cas le
		// drone vol et il retourne à son point de départ
	case Hovering:
		sd.DroneFlyingStatus.IsLanded = false
		sd.DroneFlyingStatus.IsHovering = true
		sd.DroneFlyingStatus.IsMoving = false
		sd.DroneFlyingStatus.IsPreparing = false
		sd.DroneFlyingStatus.IsGoingHome = false

	case EmergencyLanding, MotorRamping:
		sd.DroneFlyingStatus.IsLanded = false
		sd.DroneFlyingStatus.IsHovering = false
		sd.DroneFlyingStatus.IsMoving = false
		sd.DroneFlyingStatus.IsPreparing = true
		sd.DroneFlyingStatus.IsGoingHome = false
	}

	TransmitFlyingStatusUpdate(sd.DroneFlyingStatus)
	return OnSuccess
}

// UpdateFlyingState Version de OnFlyingStateReiceved sans retour
func (sd *SchedulerData) UpdateFlyingState() {
	sd.OnFlyingStateReceived()
}

// OnLastCommandSuccess La dernière commande envoyée est en succès
func (sd *SchedulerData) OnLastCommandSuccess() FuncResult {
	lastSuccess := GetLastSuccess(sd.DroneName)
	UpdateLastSuccess(sd.DroneName, NoCommand)
	log.Println("Dernière commande en succès pour ", sd.DroneName, " : ", lastSuccess)
	if sd.IntermediateCommand != NoCommand && lastSuccess != sd.IntermediateCommand {
		return OnSuccess // On passe en "Idle" directement
	}

	if sd.LastCommand != sd.IntermediateCommand && sd.IntermediateCommand != NoCommand {
		sd.IntermediateCommand = NoCommand
		return DroneReady
	}

	if sd.LastCommand == GoTo {
		return ResumeInstruction
	}

	return OnSuccess
}

// OnTakeOffEvent Décollage
func (sd *SchedulerData) OnTakeOffEvent() FuncResult {
	if sd.DroneFlyingStatus.IsLanded {
		sd.LastCommand = TakeOff
		sd.IntermediateCommand = NoCommand
		payload := &DroneCommandMessage{
			Name:   sd.LastCommand,
			Target: sd.DroneName,
		}
		TransmitEvent(payload)
		return OnSuccess
	}
	return OnError // On est déjà en vol
}

// OnGoHomeEvent Retour maison
func (sd *SchedulerData) OnGoHomeEvent() FuncResult {
	if sd.DroneFlyingStatus.IsLanded {
		sd.LastCommand = GoHome
		sd.IntermediateCommand = NoCommand
		payload := &DroneCommandMessage{
			Name:   sd.LastCommand,
			Target: sd.DroneName,
		}
		TransmitEvent(payload)
		return OnSuccess
	}
	return OnError
}
