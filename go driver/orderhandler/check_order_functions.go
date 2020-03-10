package orderhandler

import "../config"
import "../elevio"

func FindDirection(elevator *config.ElevatorState) config.Directions{
  switch elevator.Dir {
  case config.Stop:
    if checkOrdersAbove(elevator){
      return config.MovingUp
    }
    if checkOrdersBelow(elevator){
      return config.MovingDown
    }else{return config.Stop}
  case config.MovingUp:
    if checkOrdersAbove(elevator){
      return config.MovingUp
      }
    if checkOrdersBelow(elevator){
      return config.MovingDown
    } else{return config.Stop}
  case config.MovingDown:
    if checkOrdersBelow(elevator){
      return config.MovingDown
    }
    if checkOrdersAbove(elevator){
      return config.MovingUp
    } else {return config.Stop}

  }
  return config.Stop
}

func checkOrdersAbove(elevator *config.ElevatorState) bool{
  for floor := elevator.Floor + 1; floor < config.NumFloors; floor++{
    for button := 0; button < config.NumBtns; button++{
      if elevator.Queue[floor][button]{
        return true
      }
    }
  }
  return false
}

func checkOrdersBelow(elevator *config.ElevatorState) bool{
  for floor := elevator.Floor - 1; floor > -1; floor--{
    for button := 0; button < config.NumBtns; button++{
      if elevator.Queue[floor][button]{
        return true
      }
    }
  }
  return false
}


func CheckIfArrived(floor int, elevator *config.ElevatorState) bool{
  switch elevator.Dir {
  case config.MovingUp:
    if elevator.Queue[floor][elevio.BT_Cab] || elevator.Queue[floor][elevio.BT_HallUp] || !checkOrdersAbove(elevator){
      switchOffButtonLight(floor)
      clearOrderQueue(floor, elevator)
      return true
    }
  case config.MovingDown:
    if elevator.Queue[floor][elevio.BT_Cab] || elevator.Queue[floor][elevio.BT_HallDown] || !checkOrdersBelow(elevator){
      switchOffButtonLight(floor)
      clearOrderQueue(floor, elevator)
      return true
    }
  }
  return false
}


func CheckOrderSameFLoor(elevator *config.ElevatorState) bool{
  if elevator.Queue[elevator.Floor][elevio.BT_Cab] ||
      elevator.Queue[elevator.Floor][elevio.BT_HallUp] ||
        elevator.Queue[elevator.Floor][elevio.BT_HallDown]{
          switchOffButtonLight(elevator.Floor)
          clearOrderQueue(elevator.Floor, elevator)
          return true
  }
  return false
}

func clearOrderQueue(floor int, elevator *config.ElevatorState){
  elevator.Queue[floor][elevio.BT_HallUp] = false
  elevator.Queue[floor][elevio.BT_Cab] = false
  elevator.Queue[floor][elevio.BT_HallDown] = false
}

func switchOffButtonLight(floor int){
  for button := elevio.BT_HallUp; button < config.NumBtns; button++{
    elevio.SetButtonLamp(button, floor, false)
  }
}
