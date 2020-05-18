package main

import (
	"log"
	"strings"

	graphs "github.com/alonsovidales/go_graph"
)

// GetAllStreets Génère toutes les rues d'un fichier OSM décodé
func GetAllStreets() {
	for _, way := range targetMap.Ways {

		if len(way.RTags) > 0 {
			highwayTag, hok := way.GetHighwayTag()

			if hok && (!strings.Contains(highwayTag, "motorway") && !strings.Contains(highwayTag, "trunk") && !strings.Contains(highwayTag, "primary") && !strings.Contains(highwayTag, "secondary")) { //highwayTag == "crossing" || highwayTag == "living_street" || highwayTag == "unclassified" || highwayTag == "service" || highwayTag == "residential" || highwayTag == "tertiary" || highwayTag == "footway"  )  {
				name, nok := way.GetStreetName()

				if !nok {
					name = "Inconnu"
				}
				streetDataSet.Edges = append(streetDataSet.Edges, way.GenerateEdgesFromStreet(name, targetMap)...)
			}
		}
	}

	streetDataSet.Boundaries = &FlightBounds{
		MinLon: targetMap.Bounds.Minlon,
		MaxLon: targetMap.Bounds.Maxlon,
		MinLat: targetMap.Bounds.Minlat,
		MaxLat: targetMap.Bounds.Maxlat,
	}

}

// GenerateLocalID Deprecated: Génération d'un ID local (hors int64)
func (myMap *Map) GenerateLocalID() {
	for index := range myMap.Nodes {
		myMap.Nodes[index].LocalID = index + 1
	}
}

// GenerateMap Génération du graph après la section GetAllStreets()
func GenerateMap() { //
	final := make([]graphs.Edge, 0)
	log.Println("Génération de la carte")
	for index := range streetDataSet.Edges {

		final = append(final, graphs.Edge{uint64(streetDataSet.Edges[index].Start.Elem.ID), uint64(streetDataSet.Edges[index].End.Elem.ID), streetDataSet.Edges[index].Weight})
	}
	// Première itération, on a toute la carte, incluant les noeuds isolés
	streetDataSet.CityGraph = graphs.GetGraph(final, false)

}
