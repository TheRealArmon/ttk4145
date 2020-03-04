package fsm

import "../elevio"
import "../config"

func FindDirection(elevator config.ElevatorState) elevio.MotorDirection {
  switch elevator.Dir {
  case config.Stop:
    if checkOrdersAbove(elevator){
      return elevio.MD_Up
    }
    if checkOrdersBelow(elevator){
      return elevio.MD_Down
    }
  case config.MovingUp:
    if checkOrdersAbove(elevator){
      return elevio.MD_Up
      }
    if checkOrdersBelow(elevator){
      return elevio.MD_Down
    }
  case config.MovingDown:
    if checkOrdersBelow(elevator){
      return elevio.MD_Down
    }
    if checkOrdersAbove(elevator){
      return elevio.MD_Up
    }

  }
  return elevio.MD_Stop
}

func checkOrdersAbove(elevator config.ElevatorState) bool{
  for floor := 0; floor < config.NumFloors; floor++{
    for button := 0; button < config.NumBtns; button++{
      if elevator.Queue[floor][button]{
        return true
      }
    }
  }
  return false
}

func checkOrdersBelow(elevator config.ElevatorState) bool{
  for floor := config.NumFloors-1; floor > -1; floor--{
    for button := 0; button < config.NumBtns; button++{
      if elevator.Queue[floor][button]{
        return true
      }
    }
  }
  return false
}
