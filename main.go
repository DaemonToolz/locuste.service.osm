package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/keybase/go-ps"
)

var targetMap *Map
var streetDataSet FlightGraph

func main() {

	processes, err := ps.Processes()
	if err != nil {
		failOnError(err, "Error :")
	}
	procCount := 0
	for index := range processes {

		if strings.Contains(os.Args[0], processes[index].Executable()) {
			procCount++
		}

		if procCount > 1 {
			return
		}
	}

	initHealthMonitor()
	initConfiguration()
	initDroneConfiguration()
	initModuleRestartMapper()
	prepareLogs()
	pipeMain()
	initFlightSchedulerWorker()

	log.Print("Service de logging opérationnel")
	RestartRPCServer()
	initRPCClient()
	AddOrUpdateStatus(SchedulerMapHandler, false)

	for _, name := range ExtractDroneNames() {
		data := GetScheduler(name)
		TransmitAutopilotUpdate(data.Statuses) // On envoi la copie
	}

	targetMap, _ = DecodeFile("/home/pi/project/locuste/data/scan_sector.osm")
	log.Print("Carte décodée")
	targetMap.GenerateLocalID() // A mettre à jour lors de l'intégration Map-to-sectors
	streetDataSet = FlightGraph{
		MyMap: targetMap,
		Edges: make([]FlightEdge, 0),
	}
	log.Println("Graphe initialisé, traitement des informations")

	GetAllStreets()
	log.Println("Limites de la zone à scanner envoyés")

	log.Println("Carte traitée")
	for !streetDataSet.IsEulerian() {
		streetDataSet.MakeEulerian()
	}

	GenerateMap()

	log.Println("Eléments supplémentaires intégrés, on indique que les autopilotes sont prêts")

	AddOrUpdateStatus(SchedulerMapHandler, true)

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, os.Kill)

	select {
	case <-sigChan:

		AddOrUpdateStatus(SchedulerMapHandler, false)

		InterruptSchedulers()
		time.Sleep(2 * time.Second)
		StopSchedulers()
		OnServerShutdown()
		time.Sleep(1 * time.Second)
		pulse.Stop()
		go func() { stopCondition <- true }()
		Unregister()
		time.Sleep(5 * time.Second)
		logFile.Close()
		os.Exit(0)
	}
}
