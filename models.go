package main

import (
	"encoding/xml"
	"math"

	graphs "github.com/alonsovidales/go_graph"
)

/*
	Regroupement des éléments purements dédiés à la partie GO
	Marquer : https://github.com/supermitch/Chinese-Postman/
*/

// DroneIdentifier Informations envoyées à la GUI pour reconnaître le drone ciblé
type DroneIdentifier struct {
	Name string `json:"name"`
}

// FlightBounds struct
type FlightBounds struct {
	MinLat float64 `json:"min_lat"`
	MinLon float64 `json:"min_lon"`
	MaxLat float64 `json:"max_lat"`
	MaxLon float64 `json:"max_lon"`
}

// FlightCoordinate struct
type FlightCoordinate struct {
	Name string  `json:"name"`
	Lat  float64 `json:"latitude"`
	Lon  float64 `json:"longitude"`
}

// Map struct
type Map struct {
	Bounds Bounds
	Nodes  []Node
	Ways   []Way
}

// Bounds struct
type Bounds struct {
	XMLName xml.Name `xml:"bounds"`
	Minlat  float64  `xml:"minlat,attr"`
	Minlon  float64  `xml:"minlon,attr"`
	Maxlat  float64  `xml:"maxlat,attr"`
	Maxlon  float64  `xml:"maxlon,attr"`
}

// Location struct
type Location struct {
	Type        string
	Coordinates []float64
}

// Tag struct
type Tag struct {
	XMLName xml.Name `xml:"tag"`
	Key     string   `xml:"k,attr"`
	Value   string   `xml:"v,attr"`
}

// Elem Elément OSM
type Elem struct {
	ID  int64 `xml:"id,attr"`
	Loc Location
}

// Node structure
type Node struct {
	Elem
	XMLName xml.Name `xml:"node"`
	Lat     float64  `xml:"lat,attr"`
	Lng     float64  `xml:"lon,attr"`
	LocalID int
	//Tag     []Tag    `xml:"tag"`
}

// GetBoundNode Récupère un noeud à partir d'un id (int64 vs uint64)
func (mmap *Map) GetBoundNode(id int64) (*Node, bool) {
	for index := range mmap.Nodes {
		if mmap.Nodes[index].Elem.ID == id {
			return &mmap.Nodes[index], true
		}
	}
	return nil, false
}

// Way struct
type Way struct {
	Elem
	XMLName xml.Name `xml:"way"`
	RTags   []Tag    `xml:"tag"`
	Nds     []struct {
		ID int64 `xml:"ref,attr"`
	} `xml:"nd"`
}

// GetStreetName Récupère le nom d'une rue
func (way *Way) GetStreetName() (string, bool) {
	for index := range way.RTags {
		if way.RTags[index].Key == "name" {
			return way.RTags[index].Value, true
		}
	}
	return "", false
}

// GetHighwayTag Récupère le type de route
func (way *Way) GetHighwayTag() (string, bool) {
	for index := range way.RTags {
		if way.RTags[index].Key == "highway" {
			return way.RTags[index].Value, true
		}
	}
	return "", false
}

/*
	DistanceBetweenNodes Calcule la distance entre 2 points géographiques
	https://www.geodatasource.com/developers/go
*/
func DistanceBetweenNodes(nodeA Node, nodeB Node) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * nodeA.Lat / 180)
	radlat2 := float64(PI * nodeA.Lat / 180)

	theta := float64(nodeA.Lng - nodeB.Lng)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515
	dist = (dist * 1.609344) * 1000 // On veut des mètres
	return dist
}

/*
	FlightGraph Graph global englobant la carte récupérée, carte traitée et méta-données (liaisons noeud à noeud)
	Noms modifiés pour éviter la confusion avec "github.com/alonsovidales/go_graph"
*/
type FlightGraph struct {
	MyMap      *Map         `json:"map"`
	Edges      []FlightEdge `json:"edges"`
	CityGraph  *graphs.Graph
	Boundaries *FlightBounds
}

/*
	FlightEdge Pont entre 2 noeuds
	Noms modifiés pour éviter la confusion avec "github.com/alonsovidales/go_graph"
*/
type FlightEdge struct {
	End    Node    `json:"head"`
	Start  Node    `json:"tail"`
	Weight float64 `json:"weight"`
	Name   string  `json:"street_name"`
}

// GetEdge Récupère un pont / edge à partir des IDs du noeud de départ et d'arrivée
func (graph *FlightGraph) GetEdge(startId uint64, endId uint64) *FlightEdge {
	for index := range graph.Edges {
		if uint64(graph.Edges[index].Start.Elem.ID) == startId && uint64(graph.Edges[index].End.Elem.ID) == endId {
			return &graph.Edges[index]
		}
	}

	return nil
}

// GetEdges Récupère tous les ponts / edges à partir de l'ID de l'élément entrant ou sortant
func (graph *FlightGraph) GetEdges(startId uint64) *[]FlightEdge {
	edges := make([]FlightEdge, 0)
	for index := range graph.Edges {
		if uint64(graph.Edges[index].Start.Elem.ID) == startId || uint64(graph.Edges[index].End.Elem.ID) == startId {
			edges = append(edges, graph.Edges[index])
		}
	}

	return &edges
}

// InvertEdge Permet de créer une connectique inversée à partir d'un pont / edge
func (graph *FlightGraph) InvertEdge(startId uint64, endId uint64) *FlightEdge {
	var edge *FlightEdge = nil
	for index := range graph.Edges {
		if uint64(graph.Edges[index].Start.Elem.ID) == startId && uint64(graph.Edges[index].End.Elem.ID) == endId {
			edge = &graph.Edges[index]
		}
	}

	if edge != nil {
		inverted := &FlightEdge{
			Name:   edge.Name,
			End:    edge.Start,
			Start:  edge.End,
			Weight: edge.Weight,
		}

		graph.Edges = append(graph.Edges, *inverted)

		return inverted
	}
	return nil
}

// GetNode Récupère un noeud à partir de son ID
func (graph *FlightGraph) GetNode(nodeId uint64) *Node {
	for index := range graph.MyMap.Nodes {
		if uint64(graph.MyMap.Nodes[index].Elem.ID) == nodeId {
			return &graph.MyMap.Nodes[index]
		}
	}

	return nil
}

// GenerateEdgesFromStreet Permet de générer tous les ponts d'une rue donnée
func (street *Way) GenerateEdgesFromStreet(name string, myMap *Map) []FlightEdge {
	edges := make([]FlightEdge, 0)
	var head *Node
	var tail *Node

	for index, _ := range street.Nds {
		head = nil
		tail = nil

		if index >= len(street.Nds)-1 {
			continue
		}

		streetNode := street.Nds[index]
		nextStreetNode := street.Nds[index+1]

		for nodeIndex := range myMap.Nodes {
			if myMap.Nodes[nodeIndex].Elem.ID == streetNode.ID {
				tail = &myMap.Nodes[nodeIndex]
			}

			if myMap.Nodes[nodeIndex].Elem.ID == nextStreetNode.ID {
				head = &myMap.Nodes[nodeIndex]
			}
		}

		if tail != nil && head != nil {
			edges = append(edges, FlightEdge{
				Name:   name,
				End:    *head,
				Start:  *tail,
				Weight: DistanceBetweenNodes(*head, *tail),
			})
		}

	}

	return edges
}

// GetEnd Récupère le noeud opposé d'un pont
func (edge *FlightEdge) GetEnd(start Node) *Node {
	if start.LocalID == edge.Start.LocalID {
		return &edge.End
	}

	if start.LocalID == edge.End.LocalID {
		return &edge.Start
	}

	return nil
}

// OddNodes Récupère tous les noeuds ayant un nombre de connexion impaires
func (graph *FlightGraph) OddNodes() []*Node {
	oddNodes := make([]*Node, 0)
	for index := range graph.MyMap.Nodes {
		choices, _ := graph.MyMap.Nodes[index].GetChoices(graph)
		if len(*choices)%2 != 0 {
			oddNodes = append(oddNodes, &graph.MyMap.Nodes[index])
		}
	}
	return oddNodes
}

// FindDeadEnds Récupère tous les culs-de-sac
func (graph *FlightGraph) FindDeadEnds() []*Node {
	deadNodes := make([]*Node, 0)

	for _, data := range graph.OddNodes() {
		choices, _ := data.GetChoices(graph)
		if len(*choices) == 1 {
			deadNodes = append(deadNodes, data)
		}
	}

	return deadNodes
}

// IsEulerian Est un graphe Eulérien
func (graph *FlightGraph) IsEulerian() bool {
	return len(graph.OddNodes()) == 0
}

// GetChoices Récupère tous les noeuds disponibles autour d'un noeud
func (node *Node) GetChoices(graph *FlightGraph) (*[]Node, *[]FlightEdge) {
	nodes := make([]Node, 0)
	edges := make([]FlightEdge, 0)

	for index := range graph.Edges {

		if node.LocalID == graph.Edges[index].Start.LocalID {
			nodes = append(nodes, graph.Edges[index].End)
			edges = append(edges, graph.Edges[index])

		}

		if node.LocalID == graph.Edges[index].End.LocalID {
			nodes = append(nodes, graph.Edges[index].Start)
			edges = append(edges, graph.Edges[index])
		}

	}

	return &nodes, &edges
}

// FindEdges Récupère tous les bords disponibles autour d'un noeud
func (node *Node) FindEdges(graph *FlightGraph) []FlightEdge {
	edges := make([]FlightEdge, 0)
	for index := range graph.Edges {
		if node.Elem.ID == graph.Edges[index].Start.Elem.ID || node.Elem.ID == graph.Edges[index].End.Elem.ID {
			edges = append(edges, graph.Edges[index])
		}
	}
	return edges
}

// MakeEulerian Transforme un graph en graph eulérien (à corriger, algorithme invalide)
func (graph *FlightGraph) MakeEulerian() {
	deadEnds := graph.FindDeadEnds()

	for index := range deadEnds {
		edges := deadEnds[index].FindEdges(graph)

		if len(edges)%2 != 0 {
			for edgeIndex := range edges {
				newEdge := FlightEdge{
					Start:  edges[edgeIndex].End,
					End:    edges[edgeIndex].Start,
					Name:   edges[edgeIndex].Name,
					Weight: edges[edgeIndex].Weight,
				}
				graph.Edges = append(graph.Edges, newEdge)

			}
		}
	}

	oddNodes := graph.OddNodes()

	for index := range oddNodes {
		edges := oddNodes[index].FindEdges(graph)

		if len(edges)%2 != 0 {
			for edgeIndex := range edges {
				newEdge := FlightEdge{
					Start:  edges[edgeIndex].End,
					End:    edges[edgeIndex].Start,
					Name:   edges[edgeIndex].Name,
					Weight: edges[edgeIndex].Weight,
				}
				graph.Edges = append(graph.Edges, newEdge)
			}
		}
	}
}

// TotalWeight Récupère la charge totale de tous les noeuds d'un graphe
func (graph *FlightGraph) TotalWeight() float64 {
	sum := 0.0

	for index := range graph.Edges {
		sum += graph.Edges[index].Weight
	}
	return sum
}
