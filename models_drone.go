package main

// DroneCommand Commande disponible pour le drone
type DroneCommand string

// Axis Axe d'un graphe
type Axis string

const (
	// GoTo Déplacment automatique
	GoTo DroneCommand = "AutomaticGoTo"
	// Stop Annulation du déplacment automatique
	Stop DroneCommand = "AutomaticCancelGoTo"
	// TakeOff Décollage
	TakeOff DroneCommand = "AutomaticTakeOff"
	// GoHome Ordre de retour à la maison
	GoHome DroneCommand = "AutomaticGoHome"
	// Land Ordre d'atterrissage
	Land DroneCommand = "AutomaticLanding"
	// NoCommand Aucun ordre
	NoCommand DroneCommand = "NoCommand"
)

// CommandIdentifier Le "acknowledge" d'un drone pour une commande spécifique
// Envoyé par le drone
type CommandIdentifier struct {
	Name    string       `json:"name"`
	Command DroneCommand `json:"command"`
}

// DroneCommandMessage Ordre à envoyer aux drones
type DroneCommandMessage struct {
	// Name Nom de la commande
	Name DroneCommand `json:"command"`
	// Drone cible
	Target string `json:"name"`
}

// DroneFlightCoordinates Coordonnées de vol
type DroneFlightCoordinates struct {
	DroneName string            `json:"drone_name"`
	Component *FlightCoordinate `json:"coordinates"`
	Metadata  *NodeMetaData     `json:"metadata"`
}

// NodeMetaData Métadonnées de vol
type NodeMetaData struct {
	Name     string           `json:"street_name"`
	Distance float64          `json:"distance"`
	Altitude float64          `json:"altitude"`
	Previous FlightCoordinate `json:"previous"`
	Next     FlightCoordinate `json:"next"`
}

// PyDroneFlyingStatus Type d'état de vol
type PyDroneFlyingStatus int

const (
	// Landed Etat
	Landed PyDroneFlyingStatus = iota
	// TakingOff Etat
	TakingOff PyDroneFlyingStatus = iota
	// Hovering Etat
	Hovering PyDroneFlyingStatus = iota //
	// Flying Etat
	Flying PyDroneFlyingStatus = iota
	// Emergency Etat
	Emergency PyDroneFlyingStatus = iota
	// UserTakeOff Etat
	UserTakeOff PyDroneFlyingStatus = iota
	// MotorRamping Etat
	MotorRamping PyDroneFlyingStatus = iota
	// EmergencyLanding Etat
	EmergencyLanding PyDroneFlyingStatus = iota
	// None Aucun Etat
	None PyDroneFlyingStatus = iota
)
