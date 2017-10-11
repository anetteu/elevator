package elevatorManager

import (
	"../backup"
	. "../config"
	"../driver"
	"../externalOrders"
	"../network"
	"../queue"
	"../timer"
	"log"
	"os"
	"os/signal"
	"time"
)

func Iinit() {

}

func Elevator() {

	ch := NetworkCh{
		make(chan OrderMsg), make(chan OrderMsg),
		make(chan ElevatorState), make(chan ElevatorState),
		make(chan string), make(chan string),
		make(chan PeerUpdate)}

	ElevatorManager(ch)
}

func ElevatorManager(NetworkCh NetworkCh) {

	buttonCh := make(chan ButtonMsg)
	arrivedFloorCh := make(chan int)
	orderTimeoutCh := make(chan bool)
	doorTimeoutCh := make(chan bool)
	setDoorTimerCh := make(chan bool)

	network.Bcast(NetworkCh)

	go driver.ButtonPushed(buttonCh)
	go driver.ArrivedFloor(arrivedFloorCh)
	go timer.DoorTimer(doorTimeoutCh, setDoorTimerCh)
	go externalOrders.CheckOrderTimeout(orderTimeoutCh)
	go network.UpdateActivePeers()
	go motorStop(NetworkCh.OrderOut)
	go killSafe()
	go newOrderInIdle(setDoorTimerCh, NetworkCh.OrderOut)

	stateTicker := time.NewTicker(time.Millisecond * 100).C

	for {
		select {
		case call := <-buttonCh:
			buttonPressed(call, NetworkCh.OrderOut, setDoorTimerCh)
		case floor := <-arrivedFloorCh:
			driver.SetFloorIndicator(floor)
			arrivedFloor(floor, setDoorTimerCh, NetworkCh.OrderOut)
		case <-doorTimeoutCh:
			driver.SetDoorLamp(0)
			doorTimeout()
		case order := <-NetworkCh.OrderIn:
			externalOrders.OrderIn(order)
			queue.OrderIn(order, network.ID())
		case state := <-NetworkCh.StateIn:
			externalOrders.MapElevator(state)
		case <-stateTicker:
			sendElevatorState(NetworkCh.StateOut)
		case peers := <-NetworkCh.PeerIn:
			peers = network.StripPeers(peers)
			network.MapLostPeers(peers)
		case <-orderTimeoutCh:
			queue.Unassign()
			externalOrders.Reassign(NetworkCh.OrderOut)
		}
	}
}

func elevatorInitFromBackup(orderOut chan OrderMsg) {
	driver.Dead = true
	queue.Unassign()
	externalOrders.Reassign(orderOut)
	driver.SetDoorLamp(0)
	queue.Restore()
	floor := driver.FloorSignal()

	if floor != 0 {
		driver.SetMotor(DIR_DOWN)
	} else {
		driver.SetMotor(DIR_UP)
	}
}

func motorStop(orderOut chan OrderMsg) {
	for {
		_, _, timestamp, _ := driver.HwStatus()
		inactive := time.Time{}
		if timestamp != inactive && time.Since(timestamp).Seconds() > MOTOR_TIME_LIMIT {
			elevatorInitFromBackup(orderOut)
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func killSafe() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	driver.SetMotor(DIR_STOP)
	backup.Remove()
	log.Fatalf("elevator killed")
}

func doorTimeout() {
	currentFloor := driver.FloorSignal()
	driver.SetMotor(queue.NextDirection(currentFloor))
}

func arrivedFloor(floor int, setDoorTimerCh chan bool, orderOut chan OrderMsg) {
	if queue.ShouldStop(floor) {
		driver.SetMotor(0)
		driver.SetDoorLamp(1)
		setDoorTimerCh <- true
		queue.ArrivedFloor(floor)
		externalOrders.ArrivedFloor(floor, orderOut, network.ID())
	}
}

func buttonPressed(call ButtonMsg, orderOut chan OrderMsg, setDoorTimerCh chan bool) {
	order := OrderMsg{call.Floor, call.Button, "", false}
	if call.MotorDir == DIR_STOP && call.Floor == driver.FloorSignal() {
		driver.SetDoorLamp(1)
		setDoorTimerCh <- true
	} else if call.Button < N_BUTTONS-1 {
		externalButtonPressed(order, orderOut)
	} else if call.Button == N_BUTTONS-1 {
		internalButtonPressed(order)
	}
}

func internalButtonPressed(order OrderMsg) {
	queue.AddOrder(order.Floor, order.Button)
}

func externalButtonPressed(order OrderMsg, orderOut chan OrderMsg) {
	if network.ConnectionLost() != nil {
		queue.AddOrder(order.Floor, order.Button)
	} else {
		externalOrders.NewOrder(order.Floor, order.Button, orderOut)
	}
}

func sendElevatorState(stateOut chan ElevatorState) {
	motorDir, _, _, floor := driver.HwStatus()
	msg := ElevatorState{network.ID(), floor, motorDir, queue.Queue(), time.Time{}, driver.Dead}
	stateOut <- msg
}

func newOrderInIdle(setDoorTimerCh chan bool, orderOut chan OrderMsg) {
	for {
		motorDir, _, _, floor := driver.HwStatus()
		if queue.NumberOfOrders() > 0 && motorDir == DIR_STOP && !driver.DoorOpen() {
			if queue.OrderAtFloor(floor) {
				arrivedFloor(floor, setDoorTimerCh, orderOut)
			} else {
				driver.SetMotor(queue.NextDirection(floor))
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
}
