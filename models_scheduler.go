package main

import (
	"log"
	"math/rand"
	"time"
)

// Refactoring à prévoir

// SchedulerSummarizedData Informations réduites pour les autopilotes (utilisé en communication Brain <=> Scheduler)
type SchedulerSummarizedData struct {
	DroneName   string `json:"drone_name"`
	IsActive    bool   `json:"is_active"`
	IsManual    bool   `json:"is_manual"`
	IsSimulated bool   `json:"is_simulated"`
	IsRunning   bool   `json:"is_running"`
	IsReady     bool   `json:"is_ready"`
	IsBusy      bool   `json:"is_busy"`
}

// DroneSummarizedStatus Informations réduites relatif aux drones (états de vol)
type DroneSummarizedStatus struct {
	DroneName       string                    `json:"drone_name"`
	IsPreparing     bool                      `json:"is_preparing"`
	IsMoving        bool                      `json:"is_moving"`
	IsHovering      bool                      `json:"is_hovering"`
	IsLanded        bool                      `json:"is_landed"`
	IsGoingHome     bool                      `json:"is_going_home"`
	IsHomeReady     bool                      `json:"is_home_ready"`
	IsGPSFixed      bool                      `json:"is_gps_ready"`
	ReceivedStatus  PyDroneFlyingStatus       `json:"last_status"`
	ReceivedAlert   PyDroneAlertStatus        `json:"last_alert"`
	ReceivedNavHome PyDroneNavigateHomeStatus `json:"navigate_home_status"`
}

// SchedulerData Informations de l'autopîlote / planificateur de vol
type SchedulerData struct {
	DroneName            string
	OperationIndex       int
	LastCommand          DroneCommand
	IntermediateCommand  DroneCommand
	Statuses             *SchedulerSummarizedData
	DroneFlyingStatus    *DroneSummarizedStatus
	CurrentInstruction   *DroneFlightCoordinates
	Target               *Node
	SimulatedCoordinates *FlightCoordinate
	FSM                  *FlightStateMachine
	Route                *[]Node
	Distances            *[]FlightEdge
	OnUpdateChannel      *chan Event
}

// PrepareNextCoordinate On prépare le prochain point
func (currentScheduler *SchedulerData) PrepareNextCoordinate() bool {
	if !currentScheduler.Statuses.IsReady {
		return false
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	log.Print("Récupération des dernières coordonnées : ", currentScheduler.DroneName)
	if currentScheduler.OperationIndex == len(*currentScheduler.Route) /*|| (currentScheduler.OperationIndex == len(*currentScheduler.Distances))*/ {
		log.Print("Plus aucune marche de manoeuvre : ", currentScheduler.DroneName)
		return false
	}
	startingNode := (*currentScheduler.Route)[currentScheduler.OperationIndex]
	//startingEdge := (*currentScheduler.Distances)[currentScheduler.OperationIndex]
	log.Print("Création des coordonnées : ", currentScheduler.DroneName)
	if currentScheduler.CurrentInstruction == nil {
		currentScheduler.CurrentInstruction = &DroneFlightCoordinates{
			DroneName: currentScheduler.DroneName,
			Component: &FlightCoordinate{
				Lat: startingNode.Lat,
				Lon: startingNode.Lng,
			},
			Metadata: &NodeMetaData{
				//Distance: startingEdge.Weight,
				//Name: startingEdge.Name,
				Previous: FlightCoordinate{
					Lat: startingNode.Lat,
					Lon: startingNode.Lng,
				},
				Next: FlightCoordinate{
					Lat: startingNode.Lat,
					Lon: startingNode.Lng,
				},
			},
		}
	}

	if currentScheduler.OperationIndex == 0 {
		currentScheduler.CurrentInstruction.Metadata.Previous = FlightCoordinate{
			Lat: startingNode.Lat,
			Lon: startingNode.Lng,
		}
	} else {
		currentScheduler.CurrentInstruction.Metadata.Previous = FlightCoordinate{
			Lat: currentScheduler.CurrentInstruction.Component.Lat,
			Lon: currentScheduler.CurrentInstruction.Component.Lon,
		}
	}

	if currentScheduler.OperationIndex >= len(*currentScheduler.Route)-1 {
		currentScheduler.CurrentInstruction.Metadata.Next = FlightCoordinate{
			Lat: startingNode.Lat,
			Lon: startingNode.Lng,
		}
	} else {
		nextNode := (*currentScheduler.Route)[currentScheduler.OperationIndex+1]
		currentScheduler.CurrentInstruction.Metadata.Next = FlightCoordinate{
			Lat: nextNode.Lat,
			Lon: nextNode.Lng,
		}
	}

	currentScheduler.CurrentInstruction.Component.Lat = startingNode.Lat
	currentScheduler.CurrentInstruction.Component.Lon = startingNode.Lng
	//currentScheduler.CurrentInstruction.Metadata.Distance = startingEdge.Weight
	//currentScheduler.CurrentInstruction.Metadata.Name = startingEdge.Name

	if currentScheduler.Statuses.IsSimulated {
		currentScheduler.SimulatedCoordinates.Lat = startingNode.Lat
		currentScheduler.SimulatedCoordinates.Lon = startingNode.Lng

	}

	log.Print("Instruction actuelle à jour : ", currentScheduler.DroneName)
	return true
}

// ReadInitialDroneStatus Lecture des informations "position initiale"
func (currentScheduler *SchedulerData) ReadInitialDroneStatus() bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	log.Println("Réinitialisation des informations")

	currentScheduler.Statuses.IsReady = false
	currentScheduler.Route = nil
	currentScheduler.Distances = nil // Faire bosser le GC

	if !currentScheduler.Statuses.IsSimulated {

		log.Println("Acquisition des informations")

		positionMutex.Lock()
		data, ok := DroneInitialPositions[currentScheduler.DroneName]
		positionMutex.Unlock()

		if ok && currentScheduler.Target != nil {
			log.Println("Coordonnées acquises", data)
			route, edges := currentScheduler.CreateFlightRoute(data)
			currentScheduler.Route = &route
			currentScheduler.Distances = &edges
			currentScheduler.Statuses.IsReady = true
			log.Println("Autopilote prêt à utilisation", data)
			log.Println("Routes définies : ", route, edges)
		} else {
			log.Println("Aucune information n'a été récupérée", data)
			return false
		}

	} else {
		log.Println("Mode simulation détecté")

		rand.Seed(time.Now().UnixNano())

		if currentScheduler.SimulatedCoordinates == nil {
			currentScheduler.SimulatedCoordinates = &FlightCoordinate{
				Lat: targetMap.Bounds.Minlat + rand.Float64()*(targetMap.Bounds.Maxlat-targetMap.Bounds.Minlat),
				Lon: targetMap.Bounds.Minlon + rand.Float64()*(targetMap.Bounds.Maxlon-targetMap.Bounds.Minlon),
			}
		}

		if currentScheduler.Target == nil {
			return false
		}

		log.Println("Création d'une nouvelle route", currentScheduler.SimulatedCoordinates)
		route, edges := currentScheduler.CreateFlightRoute(*currentScheduler.SimulatedCoordinates)
		currentScheduler.Route = &route
		currentScheduler.Distances = &edges
		currentScheduler.Statuses.IsReady = true
		log.Println("Routes définies : ", route, edges)
		log.Println("Route créée et autopilote simulé prêt", currentScheduler.SimulatedCoordinates)
	}

	return true
}

// StopScheduler Lecture des informations "position initiale"
func (currentScheduler *SchedulerData) StopScheduler() {
	currentScheduler.Statuses.IsReady = false
	currentScheduler.Statuses.IsRunning = false
	currentScheduler.Statuses.IsBusy = false
}

// SendEvent Envoi d'une nouvelle étape / événement
func (currentScheduler *SchedulerData) SendEvent(event Event) {
	(*currentScheduler.OnUpdateChannel) <- event
}

// Stop Arrête le worker et sa machine a état
func (currentScheduler *SchedulerData) Stop() {
	currentScheduler.FSM = nil
	close(*currentScheduler.OnUpdateChannel)
	currentScheduler.OnUpdateChannel = nil
}
