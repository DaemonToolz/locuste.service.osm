package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Config Objet lié au fichier de configuration appConfig
type Config struct {
	Host    string `json:"host"`
	RPCPort int    `json:"rpc_port"`

	SchedulerPort int `json:"scheduler_port"`
}

// Drones Objet lié au fichier de configuration drone_data
type Drones struct {
	Drones   []Drone `json:"drones"`
	Altitude float64 `json:"drone_altitude"`
}

// Drone Objet drone
type Drone struct {
	IpAddress string `json:"ip_address"`
}

var appConfig Config
var logFile os.File
var drones Drones

func (cfg *Config) rpcListenUri() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.RPCPort)
}

func (cfg *Config) rpcSchedulerPort() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.SchedulerPort)
}

func initConfiguration() {
	configFile, err := os.Open("./config/appConfig.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&appConfig)
}

func initDroneConfiguration() {
	configFile, err := os.Open("/home/pi/project/locuste/config/drone_data.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&drones)
}

func prepareLogs() {
	os.MkdirAll("./logs/", 0755)

	filename := fmt.Sprintf("./logs/%s.log", filepath.Base(os.Args[0]))
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	log.SetOutput(logFile)
}
