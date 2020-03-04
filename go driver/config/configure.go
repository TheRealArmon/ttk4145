package config

import "../elevio"

const (
	NumFloors    int = 4
	StartFloor       = 0
	NumElevators     = 3
	NumBtns          = 3
)

type State int

const (
	Idle State = iota
	Moving
)

type Button int

const (
	ButtonUp Button = iota
	ButtonDown
	ButtonInside
)

type ElevatorState struct {
	ID        int
	ElevState State
	Floor     int
	Dir       elevio.MotorDirection
	Queue     [NumFloors][NumBtns]bool
}

type FSMChannels struct {
	Drv_buttons chan elevio.ButtonEvent
	Drv_floors  chan int
	Drv_stop    chan bool
	Close_door  chan bool
}

type ElevOrder struct {
	Floor             int
	Button            elevio.ButtonType
	ExecutingElevator int
	IsDone            bool
}
