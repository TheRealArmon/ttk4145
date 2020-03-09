package orderhandler

import "fmt"
import "../config"

func FindDirection(elevator *config.ElevatorState) config.Directions{
  switch elevator.Dir {
  case config.Stop:
    if checkOrdersAbove(elevator){
      fmt.Println("Above")
      return config.MovingUp
    }
    if checkOrdersBelow(elevator){
      fmt.Println("below")
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
  for floor := 0; floor < config.NumFloors; floor++{
    for button := 0; button < config.NumBtns; button++{
      if elevator.Queue[floor][button] && floor > elevator.Floor{
        return true
      }
    }
  }
  return false
}

func checkOrdersBelow(elevator *config.ElevatorState) bool{
  for floor := config.NumFloors-1; floor > -1; floor--{
    for button := 0; button < config.NumBtns; button++{
      if elevator.Queue[floor][button] && floor < elevator.Floor{
        return true
      }
    }
  }
  return false
}
