package main

import (
	"log"
)

var droneFlightRoute *[]Node

// CreateFlightRoute Création d'une nouvelle route
func (scheduler *SchedulerData) CreateFlightRoute(coordinates FlightCoordinate) ([]Node, []FlightEdge) {
	log.Println("Création du pathfinder pour le drone :", coordinates.Name)
	route := make([]Node, 0)
	visited := make([]FlightEdge, 0)
	startingPoint := DefineClosestStartingPoint(coordinates)
	log.Println("Le noeud le plus proche a été trouvé :", startingPoint)

	route = append(route, *startingPoint)

	tour, _ := streetDataSet.CityGraph.ShortestPath(uint64(startingPoint.Elem.ID), uint64(scheduler.Target.Elem.ID))
	log.Println("Le chemin le plus court : ", tour)

	for index := range tour {
		if index < len(tour)-1 {
			nextEdge := streetDataSet.GetEdge(tour[index], tour[index+1])
			if nextEdge == nil {
				nextEdge = streetDataSet.InvertEdge(tour[index+1], tour[index]) // Si l'application tierce nous retourne de nouvelles valeurs non présentes
			}

			startingPoint = streetDataSet.GetNode(tour[index])
			visited = append(visited, *nextEdge)
			route = append(route, *startingPoint)
		}
	}

	route = append(route, *scheduler.Target)
	visited = append(visited, *streetDataSet.GetEdge(uint64(scheduler.Target.Elem.ID), uint64(startingPoint.Elem.ID)))
	return route, visited

}

// DefineClosestStartingPoint Définir le point le plus proche
func DefineClosestStartingPoint(coordinates FlightCoordinate) *Node {
	log.Println("Recherche de l'emplacement le plus proche pour le drone :", coordinates.Name)
	droneNode := Node{
		Lat: coordinates.Lat,
		Lng: coordinates.Lon,
	}

	shortestDistance := 99999999999.
	temp := shortestDistance
	var closestNode *Node

	for index := range streetDataSet.Edges {
		if temp = DistanceBetweenNodes(streetDataSet.Edges[index].Start, droneNode); temp < shortestDistance {
			shortestDistance = temp
			closestNode = &streetDataSet.Edges[index].Start
		}
	}

	return closestNode
}

func defineFlightSector() {
	// TODO
	// After having found a way to create a stable
	// Eulerian graph from the map, divide
	// the map into n sectors (all independant)
}
