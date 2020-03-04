package config
import "../elevio"

const (
  NumFloors int = 4
  StartFloor    = 2
  NumElevators  = 3
  NumBtns       = 3
)

type State int

const(
  Idle State = iota
  Moving
  Init
)

type Directions int

const(
  MovingDown Directions = iota - 1
  Stop
  MovingUp
)

type ElevatorState struct{
  ID                          int
  ElevState                   State
  Floor                       int
  Dir                         Directions
  Queue [NumFloors][NumBtns]  bool
}

type FSMChannels struct {
  NewOrderToHandle chan ElevOrder
  Drv_buttons chan elevio.ButtonEvent
  Drv_floors       chan int
  Drv_stop         chan bool
  Close_door       chan bool
}

type ElevOrder struct{
  Floor               int
  Button              elevio.ButtonType
  ExecutingElevator   int
  IsDone              bool
}
