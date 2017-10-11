package main

import (
	"./driver"
	"./elevatorManager"
	"./externalOrders"
	"./queue"
)

func main() {
	driver.Init()
	queue.Init()
	externalOrders.Init()
	elevatorManager.Elevator()
}
