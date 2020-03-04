package config
import "../elevio"

const (
  NumFloors int = 4
  StartFloor    = 0
  NumElevators  = 3
  NumBtns       = 3
)

type ElevatorState struct{
  ID int
  Floor int
  Dir elevio.MotorDirection
  Queue [NumFloors][NumBtns] bool
}

type State int

const(
  Idle State = iota
  Moving     
)
