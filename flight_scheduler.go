package main

import (
	"log"
	"sync"
)

// FlightScheduler Il déterminera les ordres de vols
var FlightScheduler map[string]SchedulerData

// DroneFlyingStatus Etat de vol pour chaque drone
var DroneFlyingStatus map[string]PyDroneFlyingStatus

// DroneInitialPositions Position initiale de chaque drone
var DroneInitialPositions map[string]FlightCoordinate

// SchedulerTarget Position cible pour chaque drone
var SchedulerTarget map[string]FlightCoordinate

// LastAcknowledge Derniers "on_command_success" remontés
var LastAcknowledge map[string]DroneCommand

var accessMutex sync.Mutex
var positionMutex sync.Mutex
var targetMutex sync.Mutex
var flyingStatMutex sync.Mutex
var lastCmdMutex sync.Mutex

func initFlightSchedulerWorker() {
	log.Println("Initialisation de l'autopilote pour chaque drone")
	FlightScheduler = make(map[string]SchedulerData)
	DroneInitialPositions = make(map[string]FlightCoordinate)
	SchedulerTarget = make(map[string]FlightCoordinate)
	DroneFlyingStatus = make(map[string]PyDroneFlyingStatus)
	LastAcknowledge = make(map[string]DroneCommand)

	for _, name := range ExtractDroneNames() {
		channel := make(chan Event)
		sd := SchedulerData{
			DroneName:           name,
			LastCommand:         NoCommand,
			IntermediateCommand: NoCommand,
			Statuses: &SchedulerSummarizedData{
				DroneName:   name,
				IsReady:     false,
				IsRunning:   false,
				IsManual:    false,
				IsSimulated: false,
				IsBusy:      false,
				IsActive:    false,
			},
			DroneFlyingStatus: &DroneSummarizedStatus{
				DroneName:      name,
				IsMoving:       false,
				IsHovering:     false,
				IsLanded:       true, // Notre état par défaut, au sol
				IsGoingHome:    false,
				IsPreparing:    true,
				ReceivedStatus: None,
			},
			OperationIndex:     0,
			CurrentInstruction: nil,
			OnUpdateChannel:    &channel,

			Route:     nil,
			Distances: nil,
			FSM: &FlightStateMachine{
				Name:            name,
				DefaultFailover: OnInterrupt,
			},
		}
		sd.InitStateMachine()
		UpdateMapStatus(name, sd)
		// TODO : Lancer le service de recharge du séquenceur pour mettre à jour la machine à état
		// Sans devoir relancer tout le programme
		go StartWorker(name)
	}
}

// StartWorker Démarrage d'un autopilote
func StartWorker(name string) {
	currentScheduler := GetScheduler(name)
	// On évite de rentrer en collision avec un autre thread
	if currentScheduler.Statuses.IsRunning {
		log.Println("Unité ", currentScheduler.DroneName, "déjà en cours d'exécution")
		return
	}

	defer func(scheduler *SchedulerData) {
		if r := recover(); r != nil {
			log.Println(r)
			currentScheduler.Statuses.IsRunning = false
			UpdateMapStatus(name, *currentScheduler)
			AddOrUpdateStatus(SchedulerFlightManager, false)
		}
	}(currentScheduler)
	AddOrUpdateStatus(SchedulerFlightManager, true)
	currentScheduler.Statuses.IsRunning = true
	UpdateMapStatus(name, *currentScheduler)

	log.Println("Démarrage pour l'unité ", name)

	//var onUpdate bool = false
	//var nextStep Event = Event(-1)
SchedulerLoop:
	for currentScheduler.Statuses.IsRunning {
		//onUpdate = false
		//nextStep = Event(-1)
		select {
		case input := <-(*currentScheduler.OnUpdateChannel):
			switch input {
			case OnInterrupt: // Lui ne doit pas dans la machine a état
				currentScheduler.StopScheduler()
				UpdateMapStatus(name, *currentScheduler)
				log.Println("Arrêt de l'autopilote pour ", name)
				break SchedulerLoop

			case Idle:
				if currentScheduler.Statuses.IsBusy == true {
					currentScheduler.Statuses.IsBusy = false
					UpdateMapStatus(name, *currentScheduler)
				}

			case AskForUpdate, OnAutopilotOn, OnAutopilotOff, OnSimulation, OnNormal, SwitchedToManual, SwitchedToAutomatic:
				go currentScheduler.FSM.OnEvent(input)

			default:
				if currentScheduler.Statuses.IsActive {
					if currentScheduler.Statuses.IsBusy == false {
						currentScheduler.Statuses.IsBusy = true
						UpdateMapStatus(name, *currentScheduler)
					}
					go currentScheduler.FSM.OnEvent(input)
				}
			}

		}

	}

	currentScheduler.Statuses.IsRunning = false
	currentScheduler.Statuses.IsReady = false
	currentScheduler.Statuses.IsBusy = false
	currentScheduler.Statuses.IsActive = false

	UpdateMapStatus(name, *currentScheduler)
	AddOrUpdateStatus(SchedulerFlightManager, false) // On doit avoir 4 autopilotes, mais si un tombe HS => On remonte le tout en erreur et on applique une stratégie au cas-par-cas
}

// RestartSchedulers Redémarrer les planificateurs
func RestartSchedulers() {
	for _, name := range ExtractDroneNames() {
		go StartWorker(name)
	}
}

// UpdateInitialDroneCoordinate Mise à jour des coordonnées initiales du drone
func UpdateInitialDroneCoordinate(name string, initialCoordinates FlightCoordinate) {
	positionMutex.Lock()
	DroneInitialPositions[name] = initialCoordinates
	positionMutex.Unlock()
	go GetScheduler(name).SendEvent(NewRouteAvailable)
	log.Println("Position d'origine pour ", name, " mis à jour")
}

// UpdateMapStatus Mise à jour des status du Scheduler
func UpdateMapStatus(name string, input SchedulerData) {
	accessMutex.Lock()
	FlightScheduler[name] = input
	accessMutex.Unlock()
	TransmitAutopilotUpdate(input.Statuses) // On envoi la copie
	log.Println("Informations des status de ", name, " mis à jour")
}

// GetScheduler Récupération du planificateur
func GetScheduler(name string) *SchedulerData {
	accessMutex.Lock()
	data := FlightScheduler[name]
	accessMutex.Unlock()
	return &data

}

// InterruptSchedulers Interruption des autopilotes
func InterruptSchedulers() {
	for _, name := range ExtractDroneNames() {
		scheduler := GetScheduler(name)
		if scheduler != nil {
			scheduler.SendEvent(OnInterrupt)
		}
	}
}

// StopSchedulers Arrêt des autopilotes
func StopSchedulers() {
	for _, name := range ExtractDroneNames() {
		scheduler := GetScheduler(name)
		if scheduler != nil {
			scheduler.Stop()
		}
	}
}

// OnCommandSuccess Succès de la dernière commande
func OnCommandSuccess(identifier CommandIdentifier) {
	scheduler := GetScheduler(identifier.Name)
	if &scheduler != nil {
		UpdateLastSuccess(identifier.Name, identifier.Command)
		scheduler.SendEvent(OnCommandSuccessEvent)
		//scheduler.SendEvent(PositionReached) //
	}
}

// SendUpdateToScheduler Envoi des mises à jours
func SendUpdateToScheduler(data SchedulerSummarizedData) {
	scheduler := GetScheduler(data.DroneName)
	if scheduler != nil {
		log.Println(scheduler)
		log.Println(data)

		if scheduler.Statuses.IsActive != data.IsActive {
			if data.IsActive {
				scheduler.SendEvent(OnAutopilotOn)
			} else {
				scheduler.SendEvent(OnAutopilotOff)
			}
		}
		if scheduler.Statuses.IsManual != data.IsManual {
			if data.IsManual {
				scheduler.SendEvent(SwitchedToManual)
			} else {
				scheduler.SendEvent(SwitchedToAutomatic)
			}
		}
		if scheduler.Statuses.IsSimulated != data.IsSimulated {
			if data.IsSimulated {
				scheduler.SendEvent(OnSimulation)
			} else {
				scheduler.SendEvent(OnNormal)
			}
		}
		log.Println("Demande de mise à jour de l'autopilote pour ", data.DroneName)
	}

}

// SendEventToScheduler Envoi d'évènements
func SendEventToScheduler(data interface{}, event Event) {
	var scheduler *SchedulerData = nil

	if result, ok := data.(string); ok {
		scheduler = GetScheduler(result)
	}

	if result, ok := data.(DroneIdentifier); ok {
		scheduler = GetScheduler(result.Name)
	}

	if result, ok := data.(FlightCoordinate); ok {
		scheduler = GetScheduler(result.Name)
	}
	log.Println("Envoi de l'événement ", event, " à ", data)

	if scheduler != nil {
		scheduler.SendEvent(event)
		log.Println("Envoi de l'événement ", event, " à ", scheduler.DroneName)
	}

}

// UpdateSchedulerTarget Mise à jour de la cible (drone)
func UpdateSchedulerTarget(target FlightCoordinate) {
	targetMutex.Lock()
	SchedulerTarget[target.Name] = target
	targetMutex.Unlock()
	log.Println("Mise à jour de la cible ", target.Name)
}

// GetSchedulerTarget Récupère la cible
func GetSchedulerTarget(name string) FlightCoordinate {
	targetMutex.Lock()
	defer targetMutex.Unlock()
	return SchedulerTarget[name]

}

// UpdateFlyingStatus Mise à jour de la cible (drone)
func UpdateFlyingStatus(name string, input PyDroneFlyingStatus) {
	flyingStatMutex.Lock()
	DroneFlyingStatus[name] = input
	flyingStatMutex.Unlock()
	SendEventToScheduler(name, OnFlyingStateUpdate)
	log.Println("Mise à jour des informations de vol ", name)
}

// GetFlyingStatus Récupère la cible
func GetFlyingStatus(name string) PyDroneFlyingStatus {
	flyingStatMutex.Lock()
	defer flyingStatMutex.Unlock()
	return DroneFlyingStatus[name]

}

// UpdateLastSuccess Mise à jour de la cible (drone)
func UpdateLastSuccess(name string, input DroneCommand) {
	lastCmdMutex.Lock()
	LastAcknowledge[name] = input
	lastCmdMutex.Unlock()
	log.Println("Mise à jour des informations de vol ", name)
}

// GetLastSuccess Récupère la cible
func GetLastSuccess(name string) DroneCommand {
	lastCmdMutex.Lock()
	defer lastCmdMutex.Unlock()
	return LastAcknowledge[name]

}
