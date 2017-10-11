package config

import (
	"time"
)

type ExternalOrder struct {
	Status    int
	ID        string
	Timestamp time.Time
}

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type LostPeer struct {
	ID        string
	Timestamp time.Time
}

type ElevatorState struct {
	ID        string
	Floor     int
	Direction int
	Queue     [N_FLOORS][N_BUTTONS]int
	Timestamp time.Time
	Dead      bool
}

type OrderMsg struct {
	Floor     int
	Button    int
	ID        string
	Completed bool
}

type ButtonMsg struct {
	Floor    int
	Button   int
	MotorDir int
}

type NetworkCh struct {
	OrderOut chan OrderMsg
	OrderIn  chan OrderMsg
	StateOut chan ElevatorState
	StateIn  chan ElevatorState
	AliveOut chan string
	AliveIn  chan string
	PeerIn   chan PeerUpdate
}

const (
	N_FLOORS        = 4
	N_BUTTONS       = 3
	N_BUTTON_EXTERN = 2

	DIR_DOWN = -1
	DIR_STOP = 0
	DIR_UP   = 1

	BUTTON_UP   = 0
	BUTTON_DOWN = 1
	BUTTON_CMD  = 2

	MOTOR_SPEED = 2800

	ORDER_TIME_LIMIT            = 10
	UNCONNECTED_PEER_TIME_LIMIT = 2
	MOTOR_TIME_LIMIT            = 10
	STATE_TIME_LIMIT            = 2
	INIT_TIME_LIMIT             = 10

	BACKUP_STRING_LENGHT = N_FLOORS * N_BUTTONS * 2
)
