package main;

// DroneCommand Commande disponible pour le drone
type DroneCommand string

// Axis Axe d'un graphe
type Axis string

const (
	// GoTo Déplacment automatique
	GoTo DroneCommand = "AutomaticGoTo"
	// Stop Annulation du déplacment automatique
	Stop DroneCommand = "AutomaticCancelGoTo"
	// CamDown Rotation à 180°C de la caméra
	CamDown DroneCommand = "AutomaticSetCameraDown"
	// CamStd Remise à 0 zéro de la caméra
	CamStd DroneCommand = "AutomaticSetStandardCamera"
	// TakeOff Décollage
	TakeOff  DroneCommand= "CommonTakeOff"
	// GoHome ORdre de retour à la maison
	GoHome DroneCommand = "CommonGoHome"
)


// DroneCommandMessage Ordre à envoyer aux drones
type DroneCommandMessage struct {
	// Name Nom del a commande
	Name DroneCommand `json:"name"`
	// Params paramètres de la commande
	Params interface{}  `json:"params"`
}

// DroneFlightCoordinates Coordonnées de vol
type DroneFlightCoordinates struct {
	DroneName string `json:"drone_name"`
	Component *FlightCoordinate `json:"coordinates"`
	Metadata  *NodeMetaData `json:"metadata"`
}

// NodeMetaData Métadonnées de vol
type NodeMetaData struct{
	Name string `json:"street_name"`
	Distance float64 `json:"distance"`
	Altitude float64 `json:"altitude"`
	Previous FlightCoordinate `json:"previous"`
	Next FlightCoordinate `json:"next"`
}