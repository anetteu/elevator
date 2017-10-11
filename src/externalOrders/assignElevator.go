package externalOrders

import (
	. "../config"
	"../network"
	"../queue"
	"sort"
	"time"
)

var states = make(map[string]ElevatorState)

func MapElevator(state ElevatorState) {
	state.Timestamp = time.Now()
	states[state.ID] = state
}

func removeUnactivePeers() {
	for ID, state := range states {
		if time.Since(state.Timestamp).Seconds() > STATE_TIME_LIMIT || state.Dead {
			delete(states, ID)
		}
	}
}

func assign(floor int, button int, orderOut chan OrderMsg) {
	if orders[floor][button].Status == UNASSIGNED {
		orders[floor][button].Status = ASSIGNED
		orders[floor][button] = ExternalOrder{ASSIGNED, BestElevator(floor, button), time.Now()}
		sendOrder(orderOut, floor, button, orders[floor][button].ID, false)
	}
}

func Reassign(orderOut chan OrderMsg) {
	sort.Strings(network.ActivePeers)
	if len(network.ActivePeers) == 0 {
		for f := 0; f < N_FLOORS; f++ {
			for b := 0; b < N_BUTTON_EXTERN; b++ {
				if orders[f][b].Status != INACTIVE {
					addAssignedOrder(f, b, network.ID())
					queue.AddOrder(f, b)
				}
			}
		}
	} else if network.ID() == network.ActivePeers[0] {
		for f := 0; f < N_FLOORS; f++ {
			for b := 0; b < N_BUTTON_EXTERN; b++ {
				if orders[f][b].Status != INACTIVE {
					orders[f][b] = ExternalOrder{UNASSIGNED, "", time.Now()}
					assign(f, b, orderOut)
				}
			}
		}
	}
}

func BestElevator(floor int, button int) string {

	removeUnactivePeers()

	if len(states) == 0 {
		return network.ID()
	}

	for len(states) == 1 {
		for ID := range states {
			return ID
		}
	}

	for ID, state := range states {
		if floor == state.Floor {
			return ID
		}
	}

	freeElevators := make([]string, 0, network.N_ELEVATORS)
	for ID, state := range states {
		if state.Direction == DIR_STOP && isEmpty(state) {
			freeElevators = append(freeElevators, ID)
		}
	}
	if len(freeElevators) == 1 {
		return freeElevators[0]
	} else if len(freeElevators) > 1 {
		return shortestDistance(freeElevators, floor, button)
	}

	elevatorsInDir := make([]string, 0, network.N_ELEVATORS)
	for ID, state := range states {
		if direction(state, floor, button) {
			elevatorsInDir = append(elevatorsInDir, ID)
		}
	}
	if len(elevatorsInDir) == 1 {
		return elevatorsInDir[0]
	} else if len(elevatorsInDir) > 1 {
		return shortestDistance(elevatorsInDir, floor, button)
	}

	return network.ID()
}

func direction(state ElevatorState, floor int, button int) bool {
	if floor < state.Floor && (state.Direction == DIR_DOWN && button == BUTTON_DOWN || state.Direction == DIR_STOP && !ordersAbove(state, state.Floor)) {
		return true
	} else if floor > state.Floor && (state.Direction == DIR_UP && button == BUTTON_UP || state.Direction == DIR_STOP && !ordersBelow(state, state.Floor)) {
		return true
	}
	return false
}

func distance(state ElevatorState, floor int, button int) int {
	if floor < state.Floor {
		return state.Floor - floor
	} else if floor > state.Floor {
		return floor - state.Floor
	}
	return 0
}

func shortestDistance(elevators []string, floor int, button int) string {
	shortestDistance := 100
	elevator := ""
	sort.Strings(elevators)
	for _, ID := range elevators {
		if distance(states[ID], floor, button) == 0 {
			return string(ID)
		} else if distance(states[ID], floor, button) < shortestDistance {
			shortestDistance = distance(states[ID], floor, button)
			elevator = string(ID)
		}
	}
	return elevator
}
