package network

import (
	. "../config"
	"./bcast"
	"./localip"
	"./peers"
)

func ID() string {
	ID, _ := localip.LocalIP()
	return ID
}

func ConnectionLost() error {
	_, err := localip.LocalIP()
	return err
}

func Bcast(NetworkCh NetworkCh) {

	orderTx := make(chan OrderMsg)
	orderRx := make(chan OrderMsg)

	stateTx := make(chan ElevatorState)
	stateRx := make(chan ElevatorState)

	peerUpdateCh := make(chan PeerUpdate)
	peerTxEnable := make(chan bool)

	go bcast.Transmitter(16560, orderTx)
	go bcast.Receiver(16560, orderRx)

	go bcast.Transmitter(16570, stateTx)
	go bcast.Receiver(16570, stateRx)

	go peers.Transmitter(15630, aliveElevators(), peerTxEnable)
	go peers.Receiver(15630, peerUpdateCh)

	go transmitOrder(NetworkCh.OrderOut, orderTx)
	go receiveOrder(NetworkCh.OrderIn, orderRx)

	go transmitState(NetworkCh.StateOut, stateTx)
	go receiveState(NetworkCh.StateIn, stateRx)

	go receivePeerUpdateCh(NetworkCh.PeerIn, peerUpdateCh)

}
func transmitOrder(orderOut chan OrderMsg, orderTx chan OrderMsg) {
	for {
		msg := <-orderOut
		for i := 0; i < 20; i++ {
			orderTx <- msg
		}
	}
}

func receiveOrder(orderIn chan OrderMsg, orderRx chan OrderMsg) {
	for {
		msg := <-orderRx
		orderIn <- msg
	}
}

func transmitState(stateOut chan ElevatorState, stateTx chan ElevatorState) {
	for {
		msg := <-stateOut
		stateTx <- msg
	}
}

func receiveState(stateIn chan ElevatorState, stateRx chan ElevatorState) {
	for {
		msg := <-stateRx
		stateIn <- msg
	}
}

func receivePeerUpdateCh(peerIn chan PeerUpdate, peerUpdateCh chan PeerUpdate) {
	for {
		msg := <-peerUpdateCh
		peerIn <- msg
	}
}
