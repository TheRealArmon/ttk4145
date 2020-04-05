package orderhandler

import "../config"
import "../elevio"
//import "fmt"
//input elevatorStatus

func FindDirection(elevatorState *config.ElevatorState) config.Directions{
  switch elevatorState.Dir {
  case config.Stop:
    if checkOrdersAbove(elevatorState){
      return config.MovingUp
    }
    if checkOrdersBelow(elevatorState){
      return config.MovingDown
    }else{return config.Stop}
  case config.MovingUp:
    if checkOrdersAbove(elevatorState){
      return config.MovingUp
      }
    if checkOrdersBelow(elevatorState){
      return config.MovingDown
    } else{return config.Stop}
  case config.MovingDown:
    if checkOrdersBelow(elevatorState){
      return config.MovingDown
    }
    if checkOrdersAbove(elevatorState){
      return config.MovingUp
    } else {return config.Stop}

  }
  return config.Stop
}

func checkOrdersAbove(elevatorState *config.ElevatorState) bool{
  for floor := elevatorState.Floor + 1; floor < config.NumFloors; floor++{
    for button := 0; button < config.NumBtns; button++{
      if elevatorState.Queue[floor][button]{
        return true
      }
    }
  }
  return false
}

func checkOrdersBelow(elevatorState *config.ElevatorState) bool{
  for floor := elevatorState.Floor - 1; floor > -1; floor--{
    for button := 0; button < config.NumBtns; button++{
      if elevatorState.Queue[floor][button]{
        return true
      }
    }
  }
  return false
}


func CheckIfArrived(floor int, elevatorState *config.ElevatorState) bool{
  switch elevatorState.Dir {
  case config.MovingUp:
    if elevatorState.Queue[floor][elevio.BT_Cab] || elevatorState.Queue[floor][elevio.BT_HallUp] || !checkOrdersAbove(elevatorState){
      ClearOrderQueue(floor, elevatorState)
      return true
    }
  case config.MovingDown:
    if elevatorState.Queue[floor][elevio.BT_Cab] || elevatorState.Queue[floor][elevio.BT_HallDown] || !checkOrdersBelow(elevatorState){
      ClearOrderQueue(floor, elevatorState)
      return true
    }
  }
  return false
}


func CheckOrderSameFLoor(elevatorState *config.ElevatorState) bool{
  floor := elevatorState.Floor
  if elevatorState.Queue[floor][elevio.BT_Cab] ||
      elevatorState.Queue[floor][elevio.BT_HallUp] ||
        elevatorState.Queue[floor][elevio.BT_HallDown]{
          ClearOrderQueue(floor, elevatorState)
          return true
  }
  return false
}

func ClearOrderQueue(floor int, elevatorState *config.ElevatorState){
  elevatorState.Queue[floor][elevio.BT_HallUp] = false
  elevatorState.Queue[floor][elevio.BT_Cab] = false
  elevatorState.Queue[floor][elevio.BT_HallDown] = false
}

func SwitchOffButtonLight(floor int){
  for button := elevio.BT_HallUp; button < config.NumBtns; button++{
    elevio.SetButtonLamp(button, floor, false)
  }
}
