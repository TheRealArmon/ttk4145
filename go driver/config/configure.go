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
  Idle State = 0
  Moving     = 1
  Init       = 2
)

type Directions int

const(
  MovingDown Directions = -1
  Stop                  = 0
  MovingUp              = 1
)

type ElevatorState struct{
  ID                          int
  Floor                       int
  Dir                         Directions
  ElevState                   State
  Queue [NumFloors][NumBtns]  bool
}

type FSMChannels struct {
  NewOrderToHandle chan ElevatorOrder
  Drv_buttons chan elevio.ButtonEvent
  Drv_floors       chan int
  Drv_stop         chan bool
  Close_door       chan bool
}

type ElevatorOrder struct{
  Button              elevio.ButtonType
  Floor               int
  ExecutingElevator   int
  OrderDone           bool
}
