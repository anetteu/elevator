package queue

import (
	"../backup"
	. "../config"
	"../driver"
	"fmt"
)

var orders = [N_FLOORS][N_BUTTONS]int{}

func Init() {
	Restore()
}

func Queue() [N_FLOORS][N_BUTTONS]int {
	return orders
}

func Restore() {
	if backup.Exists() {
		queue := backup.Read(BACKUP_STRING_LENGHT)
		orders = backup.Queue(queue)
		fmt.Println(orders)
		updateLights()
	}
}

func AddOrder(floor int, button int) {
	orders[floor][button] = 1
	driver.SetButtonLamp(floor, button, 1)
	backup.Write(backup.String(orders))
}

func DeleteOrder(floor int, button int) {
	orders[floor][button] = 0
	driver.SetButtonLamp(floor, button, 0)
	backup.Write(backup.String(orders))
}

func ArrivedFloor(floor int) {
	for b := 0; b < N_BUTTONS; b++ {
		DeleteOrder(floor, b)
	}
}

func OrderIn(order OrderMsg, ID string) {
	if order.Completed && orders[order.Floor][order.Button] == 1 {
		DeleteOrder(order.Floor, order.Button)
	} else if order.ID == ID && orders[order.Floor][order.Button] == 0 {
		AddOrder(order.Floor, order.Button)
	}
}

func Unassign() {
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTON_EXTERN; b++ {
			DeleteOrder(f, b)
		}
	}
}

func OrderAtFloor(floor int) bool {
	if floor == -1 {
		return false
	}
	for b := 0; b < N_BUTTONS; b++ {
		if orders[floor][b] == 1 {
			return true
		}
	}
	return false
}

func ShouldStop(floor int) bool {
	motorDir, _, _, _ := driver.HwStatus()
	if floor == 0 || floor == N_FLOORS-1 {
		return true
	} else if motorDir == DIR_UP && !ordersAbove(floor) && !OrderAtFloor(floor) {
		return true
	} else if motorDir == DIR_DOWN && !ordersBelow(floor) && !OrderAtFloor(floor) {
		return true
	} else if !OrderAtFloor(floor) {
		return false
	} else if orders[floor][BUTTON_CMD] == 1 || (orders[floor][BUTTON_UP] == 1 && (motorDir == DIR_UP)) || (orders[floor][BUTTON_DOWN] == 1 && (motorDir == DIR_DOWN)) {
		return true
	} else if (orders[floor][BUTTON_DOWN] == 1 && ordersAbove(floor) && motorDir == DIR_UP) || (orders[floor][BUTTON_UP] == 1 && ordersBelow(floor) && motorDir == DIR_DOWN) {
		return false
	} else {
		return true
	}
}

func NextDirection(floor int) int {
	motorDir, prevMotorDir, _, _ := driver.HwStatus()
	if motorDir == DIR_UP {
		if ordersAbove(floor) {
			return DIR_UP
		} else if ordersBelow(floor) {
			return DIR_DOWN
		}
	}
	if motorDir == DIR_DOWN {
		if ordersBelow(floor) {
			return DIR_DOWN
		} else if ordersAbove(floor) {
			return DIR_UP
		}
	}
	if motorDir == DIR_STOP {
		if prevMotorDir == DIR_UP && ordersAbove(floor) {
			return DIR_UP
		} else if prevMotorDir == DIR_DOWN && ordersBelow(floor) {
			return DIR_DOWN
		} else if ordersAbove(floor) {
			return DIR_UP
		} else if ordersBelow(floor) {
			return DIR_DOWN
		}
	}
	return DIR_STOP
}

func ordersAbove(floor int) bool {
	if floor == N_FLOORS-1 {
		return false
	}
	for f := floor + 1; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if orders[f][b] == 1 {
				return true
			}
		}
	}
	return false
}

func ordersBelow(floor int) bool {
	if floor == 0 {
		return false
	}
	for f := floor - 1; f >= 0; f-- {
		for b := 0; b < N_BUTTONS; b++ {
			if orders[f][b] == 1 {
				return true
			}
		}
	}
	return false
}

func NumberOfOrders() int {
	numberOfOrders := 0
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			numberOfOrders += orders[f][b]
		}
	}
	return numberOfOrders
}

func updateLights() {
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			driver.SetButtonLamp(f, b, orders[f][b])
		}
	}
}
