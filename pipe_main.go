package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var ongoingDiagProcess bool

// commandAssociation On bind les informations reçues avec
var commandAssociation map[string](map[string]interface{}) // On ajouter le boolean pour indiquer si on a un string param

func initCommandAssociation() {
	commandAssociation = make(map[string](map[string]interface{}))
	commandAssociation["list"] = map[string]interface{}{
		"func":         listModules,
		"string_param": false,
	}
	commandAssociation["start"] = map[string]interface{}{
		"func":         moduleStart,
		"string_param": true,
	}

	commandAssociation["restart"] = map[string]interface{}{
		"func":         moduleRestart,
		"string_param": true,
	}
	commandAssociation["stop"] = map[string]interface{}{
		"func":         moduleStop,
		"string_param": true,
	}

}

var readPipe *os.File
var writePipe *os.File
var dataMutex sync.Mutex
var sharedData []string

var in string
var out string

// AddPipeSharedData Permet d'ajouter une information en provenance de la pipe entrante
func AddPipeSharedData(input string) {
	dataMutex.Lock()
	sharedData = append(sharedData, input)
	dataMutex.Unlock()
}

// GetFirstPipeSharedData Récupère la plus vieille instruction disponible
func GetFirstPipeSharedData() string {
	var firstInstruction string
	dataMutex.Lock()
	firstInstruction, sharedData = sharedData[0], sharedData[1:]
	dataMutex.Unlock()
	return firstInstruction
}

// Diagnostic code
func pipeMain() {
	out = "/tmp/locuste.diagnostic.scheduler"
	in = "/tmp/locuste.scheduler.diagnostic"
	log.Println("Ouverture de la pipe nommée pour locuste.service.brain")
	syscall.Mkfifo(in, 0666)
	syscall.Mkfifo(out, 0666)
	initCommandAssociation()
	ongoingDiagProcess = true
	log.Println("En attente d'instruction de diagnostiques")
	go startReadProcess()
	go startWriteProcess()

}

func startReadProcess() {
	for ongoingDiagProcess == true {
		var buffer bytes.Buffer
		readPipe, _ = os.OpenFile(in, os.O_RDONLY, os.ModeNamedPipe)
		io.Copy(&buffer, readPipe)
		if buffer.Len() > 0 {
			AddPipeSharedData(buffer.String())
		}
		readPipe.Close()
		time.Sleep(250 * time.Millisecond)
	}
}

func startWriteProcess() {
	for ongoingDiagProcess == true {
		writePipe, _ := os.OpenFile(out, os.O_WRONLY, os.ModeNamedPipe)
		if len(sharedData) > 0 {
			executeDiagnosticFunction(writePipe, GetFirstPipeSharedData())
		}
		writePipe.Close()
		time.Sleep(300 * time.Millisecond)
	}
}

func executeDiagnosticFunction(output *os.File, input string) {
	log.Println("Instruction de diagnostique : ", input)
	if strings.Count(input, " ") >= 1 {
		dataStr := strings.Split(input, " ")

		if funcMapper, ok := commandAssociation[dataStr[1]]; ok {
			hasStrParam := funcMapper["string_param"].(bool)
			if hasStrParam && len(dataStr) > 2 {
				funcMapper["func"].(func(*os.File, string))(output, dataStr[2])
			} else if !hasStrParam {
				funcMapper["func"].(func(*os.File))(output)
			}
		}

	}
	log.Println("Réponse envoyée")

}

func listModules(output *os.File) {
	outputStr := ""
	for key, value := range GlobalStatuses {
		outputStr += fmt.Sprintf("[%s]:%s,", string(key), strconv.FormatBool(value))
	}

	output.WriteString(outputStr)
	log.Println("Informations envoyées " + outputStr)
}

func moduleStart(output *os.File, input string) {
	outStr := ""
	if CallModuleRestart(Component(input)) {
		outStr = fmt.Sprintf("[%s] - Module démarré avec succès", input)
	} else {
		outStr = fmt.Sprintf("[%s] - Echec du démarrage du module", input)
	}
	output.WriteString(outStr)
}

func moduleRestart(output *os.File, input string) {
	outStr := ""
	if CallModuleRestart(Component(input)) {
		outStr = fmt.Sprintf("[%s] - Module redémarré avec succès", input)
	} else {
		outStr = fmt.Sprintf("[%s] - Echec du redémarrage du module", input)
	}
	output.WriteString(outStr)
}

func moduleStop(output *os.File, input string) {

}
