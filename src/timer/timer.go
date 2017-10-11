package timer

import (
	"time"
)

func DoorTimer(doorTimeout chan bool, setDoorTimerCh chan bool) {
	const doorOpenTime = 3 * time.Second
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-setDoorTimerCh:
			timer.Reset(doorOpenTime)
		case <-timer.C:
			timer.Stop()
			doorTimeout <- true
		}
	}
}
