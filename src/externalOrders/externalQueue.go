package externalOrders

import (
	. "../config"
	"../driver"
	"time"
)

const (
	INACTIVE   = 0
	UNASSIGNED = 1
	ASSIGNED   = 2
)

var orders = [N_FLOORS][N_BUTTON_EXTERN]ExternalOrder{}

func Init() {
	clearOrders()
}

func DeleteOrder(floor int, button int) {
	orders[floor][button] = ExternalOrder{INACTIVE, "", time.Time{}}
	driver.SetButtonLamp(floor, button, 0)
}

func addAssignedOrder(floor int, button int, ID string) {
	orders[floor][button] = ExternalOrder{ASSIGNED, ID, time.Now()}
	driver.SetButtonLamp(floor, button, 1)
}

func addUnassignedOrder(floor int, button int) {
	orders[floor][button] = ExternalOrder{UNASSIGNED, "", time.Now()}
	driver.SetButtonLamp(floor, button, 1)
}

func NewOrder(floor int, button int, orderOut chan OrderMsg) {
	if orders[floor][button].Status == INACTIVE {
		addUnassignedOrder(floor, button)
		assign(floor, button, orderOut)
	}
}

func ArrivedFloor(floor int, orderOut chan OrderMsg, ID string) {
	for b := 0; b < N_BUTTON_EXTERN; b++ {
		if orders[floor][b].ID == ID {
			DeleteOrder(floor, b)
			sendOrder(orderOut, floor, b, "", true)
		}
	}
}

func OrderIn(order OrderMsg) {
	if order.Completed && orders[order.Floor][order.Button].Status != INACTIVE {
		DeleteOrder(order.Floor, order.Button)
	} else if !order.Completed {
		addAssignedOrder(order.Floor, order.Button, order.ID)
	}
}

func CheckOrderTimeout(orderTimeout chan bool) {
	for {
		for f := 0; f < N_FLOORS; f++ {
			for b := 0; b < N_BUTTON_EXTERN; b++ {
				if time.Since(orders[f][b].Timestamp).Seconds() < 100000 && time.Since(orders[f][b].Timestamp).Seconds() > ORDER_TIME_LIMIT {
					resetAllOrderTimers()
					orderTimeout <- true
				}
			}
		}
		time.Sleep(time.Second * 3)
	}
}

func sendOrder(orderOut chan OrderMsg, floor int, button int, ID string, completed bool) {
	msg := OrderMsg{floor, button, ID, completed}
	orderOut <- msg
}

func clearOrders() {
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTON_EXTERN; b++ {
			DeleteOrder(f, b)
		}
	}
}

func resetAllOrderTimers() {
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTON_EXTERN; b++ {
			orders[f][b].Timestamp = time.Time{}
		}
	}
}

func isEmpty(state ElevatorState) bool {
	queue := state.Queue
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if queue[f][b] == 1 {
				return false
			}
		}
	}
	return true
}

func ordersAbove(state ElevatorState, floor int) bool {
	var queue [N_FLOORS][N_BUTTONS]int
	queue = state.Queue
	if floor == N_FLOORS-1 {
		return false
	}
	for f := floor + 1; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if queue[f][b] == 1 {
				return true
			}
		}
	}
	return false
}

func ordersBelow(state ElevatorState, floor int) bool {
	var queue [N_FLOORS][N_BUTTONS]int
	queue = state.Queue
	if floor == 0 {
		return false
	}
	for f := floor - 1; f >= 0; f-- {
		for b := 0; b < N_BUTTONS; b++ {
			if queue[f][b] == 1 {
				return true
			}
		}
	}
	return false
}
