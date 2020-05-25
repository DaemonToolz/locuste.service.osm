package main

import (
	"fmt"
	"log"
)

// FuncResult Résultat d'un événement de séquenceur
type FuncResult int

const (
	// OnError Résultat d'erreur
	OnError FuncResult = iota
	// OnSuccess Résultat OnSuccess
	OnSuccess FuncResult = iota
	// UpdateRequired On demande une mise à jour
	UpdateRequired FuncResult = iota
	// OnTerminated En cas de fin du circuit
	OnTerminated FuncResult = iota
	// TakeOffRequired Le drone est au sol
	TakeOffRequired FuncResult = iota
	// LandingRequired Le drone doit atterrir
	LandingRequired FuncResult = iota
	// MoveInterruptRequired Le drone doit s'arrêter
	MoveInterruptRequired FuncResult = iota
	// StatusCheckRequired On va faire un check-up pour connaitre l'état de vol
	StatusCheckRequired FuncResult = iota
	// ResumeInstruction On reprend la commande d'origine
	ResumeInstruction FuncResult = iota
	// DroneReady Le drone est prêt et on peut procéder
	DroneReady FuncResult = iota
)

// Event Type d'évènement
type Event int

const (
	// AnyEvent Tout événement est autorisé
	AnyEvent Event = iota
	// Idle On ne fait rien
	Idle Event = iota
	// StartingPointDefined Le drone a déterminé un nouveau point de décollage
	StartingPointDefined Event = iota
	// NewRouteAvailable Une nouvelle route est affectée
	NewRouteAvailable Event = iota
	// PositionReached Le drone a atteint sa position
	PositionReached Event = iota
	// SwitchedToManual Bascule en pilotage manuel
	SwitchedToManual Event = iota
	// SwitchedToAutomatic Bascule en mode autopilote
	SwitchedToAutomatic Event = iota
	// OnSimulation Bascule en mode simulation
	OnSimulation Event = iota
	// OnNormal Bascule en mode normal (après le mode simulation)
	OnNormal Event = iota
	// OnInterrupt Interruption du Scheduler
	OnInterrupt Event = iota
	// OnTargetDefined Cible de navigation définiée
	OnTargetDefined Event = iota
	// OnTargetReady Cible prête
	OnTargetReady Event = iota
	// OnDestinationReached La cible a été atteinte
	OnDestinationReached Event = iota
	// AskForUpdate Demande de mise à jour (état interne)
	AskForUpdate Event = iota
	// OnAutopilotOn On active le pilote automatique
	OnAutopilotOn Event = iota
	// OnAutopilotOff On désactive le pilote automatique
	OnAutopilotOff Event = iota
	// OnDroneStatusCheckup On demande le checkup
	OnDroneStatusCheckup Event = iota
	// OnTakeOff On demande le décollage
	OnTakeOff Event = iota
	// OnLanding On demande le l'atterrissage
	OnLanding Event = iota
	// OnGoHome On demande le retour maison
	OnGoHome Event = iota
	// OnPreparing On est en cours de prépration
	OnPreparing Event = iota
	// OnFlyingStateUpdate Mise à jour du dernier état de vol
	OnFlyingStateUpdate Event = iota
	// OnCommandSuccessEvent On indique le succès du dernier événement
	OnCommandSuccessEvent Event = iota
	// IntermerdiateOrderRequired  Ordre intermédiaire requis avant de procéder
	IntermerdiateOrderRequired Event = iota
)

// FlightStateMachine Machine à état du pilote automatique
type FlightStateMachine struct {
	Name            string
	DefaultFailover Event
	currentState    State
	callbacks       map[Event]State
	eventChannel    *chan Event
}

// Init Initialisation
func (fsm *FlightStateMachine) Init(output *chan Event) {
	fsm.callbacks = make(map[Event]State)
	fsm.eventChannel = output
}

// AddState Ajout d'un état
func (fsm *FlightStateMachine) AddState(state State) {
	log.Println(fmt.Sprintf("[%s] : ", fsm.Name), "Ajout d'un événement", state)
	fsm.callbacks[state.Event] = state
}

// AddEventState Ajout d'un état
func (fsm *FlightStateMachine) AddEventState(evt Event, state State) {
	log.Println(fmt.Sprintf("[%s] : ", fsm.Name), "Ajout d'un événement ", evt, "avec l'état", state)
	fsm.callbacks[evt] = state
}

// SetStateOutcome Ajout d'une transition
func (fsm *FlightStateMachine) SetStateOutcome(input Event, onResult FuncResult, output Event) {
	log.Println(fmt.Sprintf("[%s] : ", fsm.Name), input, "+", onResult, "=", output)
	cb := fsm.callbacks[input]
	cb.SetOutcome(onResult, output)
	fsm.callbacks[input] = cb
}

// OnEvent Appel de la séquence Pré-call, call, et post-call
func (fsm *FlightStateMachine) OnEvent(evt Event) {
	if evt == Idle {
		return // Il s'agit de l'événement ignoré
	}

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			(*fsm.eventChannel) <- OnInterrupt
		}
	}()

	log.Println(fmt.Sprintf("[SCHEDULER - %s] : ", fsm.Name), "Appel de l'événement", evt)
	nextState, resultOk := fsm.callbacks[evt]
	log.Println(fmt.Sprintf("[SCHEDULER - %s] : ", fsm.Name), fsm.currentState, "--->", nextState, resultOk)

	if nextState.OperationAuthorized(fsm.currentState.Event) {
		fsm.currentState = nextState
		log.Println(fmt.Sprintf("[SCHEDULER - %s] : ", fsm.Name), "Action autorisée")
	} else {
		log.Println(fmt.Sprintf("[SCHEDULER - %s] : ERREUR - ", fsm.Name), "Action non autorisée, bascule en mode erreur")
		resultOk = false
	}

	if !resultOk {
		log.Println(fmt.Sprintf("[SCHEDULER - %s] : ERREUR - ", fsm.Name), "Echec de l'appel, fail over par défaut appelé")
		(*fsm.eventChannel) <- fsm.DefaultFailover
	} else {
		(*fsm.eventChannel) <- fsm.currentState.Call()
	}
}

// State Etat de la machine à états
type State struct {
	Name           string
	Description    string
	Event          Event
	Result         FuncResult // On focus le succès / échec, car on ne veut pas complexifier le code qui va arriver ici
	possibleStates map[FuncResult]Event
	possibleIncome []Event

	function interface{}
	precall  interface{}
	postcall interface{}
}

// Init Initialisation
func (state *State) Init() {
	state.possibleStates = make(map[FuncResult]Event)
	state.possibleIncome = make([]Event, 0)
}

// SetCallBacks Ajout d'un callback
func (state *State) SetCallBacks(precall interface{}, function interface{}, postcall interface{}) {
	if function == nil {
		panic("Impossible d'avoir un callback NULL")
	}

	state.precall = precall
	state.function = function
	state.postcall = postcall
}

// SetOutcome Ajout d'une transition résultat => état
func (state *State) SetOutcome(result FuncResult, outcome Event) {
	state.possibleStates[result] = outcome
}

// DeleteOutcome Retrait d'une transition
func (state *State) DeleteOutcome(result FuncResult) {
	delete(state.possibleStates, result)
}

// OperationAuthorized Opération invoquée autorisée
func (state *State) OperationAuthorized(income Event) bool {
	if state.possibleIncome == nil {
		return false
	}

	if len(state.possibleIncome) > 0 {
		for _, event := range state.possibleIncome {
			if event == AnyEvent || event == income {
				return true
			}
		}
		return false
	}
	return true
}

// SetIncome Ajoute un état aux états précédents autorisée
func (state *State) SetIncome(income Event) {
	state.possibleIncome = append(state.possibleIncome, income)
}

// DeleteIncome Retrait d'un état autorisé
func (state *State) DeleteIncome(income Event) {
	index := -1
	for eventIndex := range state.possibleIncome {
		if income == state.possibleIncome[eventIndex] {
			index = eventIndex
			break
		}
	}

	if index != -1 {
		var temp []Event
		temp = append(state.possibleIncome[:index], state.possibleIncome[index+1:]...)
		state.possibleIncome = temp
	}
}

// Call Appelle la séquence
func (state *State) Call() Event {
	if state.precall != nil {
		state.precall.(func())()
	}

	state.Result = state.function.(func() FuncResult)()

	if state.postcall != nil {
		state.postcall.(func())()
	}

	if len(state.possibleStates) == 0 {
		return Idle
	}

	return state.possibleStates[state.Result]

}

// StateMachineJSON Modèle Json stocké dans config/state_machine
type StateMachineJSON struct {
	StateMachine []StateJSON `json:"state_machine"`
}

// StateJSON Représentation de l'état au format JSON
type StateJSON struct {
	Name        string           `json:"name"`
	Event       int              `json:"on_event"`
	Description string           `json:"description"`
	Outcomes    []TransitionJSON `json:"outcomes"`
	Incomes     []int            `json:"incomes"`
	Callbacks   CallbackJSON     `json:"callbacks"`
}

// TransitionJSON Transition état à état en JSON
type TransitionJSON struct {
	Result    int `json:"result"`
	NextState int `json:"next_state"`
}

// CallbackJSON Callbacks en JSON
type CallbackJSON struct {
	Precall  string `json:"precall"`
	Call     string `json:"call"`
	Postcall string `json:"postcall"`
}
