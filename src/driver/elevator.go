package driver

import (
	. "../config"
	"time"
)

var motorDir int
var prevMotorDir int
var motorTimestamp time.Time
var currentFloor int
var Dead bool

var lampChannels = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}}

var buttonChannels = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4}}

func Init() {
	ioInit()
	TurnOffLamps()
	SetDoorLamp(0)
	timestamp := time.Now()
	for {
		if FloorSignal() != 0 && time.Since(timestamp).Seconds() < INIT_TIME_LIMIT {
			SetMotor(DIR_DOWN)
		} else if FloorSignal() != -1 && time.Since(timestamp).Seconds() > INIT_TIME_LIMIT {
			SetMotor(DIR_STOP)
			break
		} else if FloorSignal() != 0 {
			SetMotor(DIR_UP)
		} else {
			SetMotor(DIR_STOP)
			break
		}
	}
	Dead = false
}

func SetMotor(dir int) {
	if dir == 0 {
		prevMotorDir = motorDir
		motorDir = DIR_STOP
		writeAnalog(MOTOR, 0)
	} else if dir > 0 {
		prevMotorDir = motorDir
		motorDir = DIR_UP
		if prevMotorDir == DIR_STOP {
			motorTimestamp = time.Now()
		}
		clearBit(MOTORDIR)
		writeAnalog(MOTOR, MOTOR_SPEED)
	} else if dir < 0 {
		prevMotorDir = motorDir
		motorDir = DIR_DOWN
		if prevMotorDir == DIR_STOP {
			motorTimestamp = time.Now()
		}
		setBit(MOTORDIR)
		writeAnalog(MOTOR, MOTOR_SPEED)
	}
}

func HwStatus() (int, int, time.Time, int) {
	return motorDir, prevMotorDir, motorTimestamp, currentFloor
}

func SetButtonLamp(floor int, button int, value int) {
	if value == 1 {
		setBit(lampChannels[floor][button])
	} else {
		clearBit(lampChannels[floor][button])
	}
}

func ButtonPushed(buttonCh chan ButtonMsg) {
	for {
		for f := 0; f < N_FLOORS; f++ {
			for b := 0; b < N_BUTTONS; b++ {
				if ButtonSignal(f, b) {
					motorDir, _, _, _ := HwStatus()
					msg := ButtonMsg{f, b, motorDir}
					time.Sleep(time.Millisecond * 50)
					buttonCh <- msg
				}
			}
		}
	}
}

func ArrivedFloor(floorCh chan int) {
	floor := -1
	for {
		if FloorSignal() != floor && FloorSignal() != -1 {
			Dead = false
			floor = FloorSignal()
			currentFloor = floor
			motorTimestamp = time.Time{}
			floorCh <- floor
		}
	}
}

func SetFloorIndicator(floor int) {
	if floor&0x02 > 0 {
		setBit(LIGHT_FLOOR_IND1)
	} else {
		clearBit(LIGHT_FLOOR_IND1)
	}
	if floor&0x01 > 0 {
		setBit(LIGHT_FLOOR_IND2)
	} else {
		clearBit(LIGHT_FLOOR_IND2)
	}
}

func SetDoorLamp(value int) {
	if value == 1 {
		setBit(LIGHT_DOOR_OPEN)
	} else {
		clearBit(LIGHT_DOOR_OPEN)
	}
}

func TurnOffLamps() {
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			clearBit(lampChannels[f][b])
		}
	}
}

func ButtonSignal(floor int, button int) bool {
	return readBit(buttonChannels[floor][button])
}

func FloorSignal() int {
	if readBit(SENSOR_FLOOR1) {
		return 0
	} else if readBit(SENSOR_FLOOR2) {
		return 1
	} else if readBit(SENSOR_FLOOR3) {
		return 2
	} else if readBit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func DoorOpen() bool {
	return readBit(LIGHT_DOOR_OPEN)
}
