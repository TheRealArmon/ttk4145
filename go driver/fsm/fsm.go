package fsm

import "../elevio"
import "../config"
import "../timer"
import "../orderhandler"
//import "fmt"
//import "time"
//import "reflect"


func initState(elevator *config.ElevatorState) {
  elevio.SetDoorOpenLamp(false)
  for i := 0; i < config.NumFloors; i++{
    for j := elevio.BT_HallUp; j < config.NumBtns; j++{
      elevio.SetButtonLamp(j, i, false)
      elevator.Queue[i][j] = false
    }
  }
  elevio.SetMotorDirection(elevio.MD_Down)
}

func reachedFloor(sender <-chan bool, id string) {
  elevio.SetMotorDirection(elevio.MD_Stop)
  var status = config.ElevatorMap[id]
  status.ElevState = config.Idle
  status.Dir = config.Stop
  config.ElevatorMap[id] = status
  elevio.SetDoorOpenLamp(true)
  for {
    select{
    case <- sender:
      elevio.SetDoorOpenLamp(false)
      return
    }
  }
}



func ElevStateMachine(ch config.FSMChannels, newOrder chan config.ElevatorOrder, id string, reciever chan<- map[string]config.ElevatorState) {

  elevator := config.ElevatorState{
    Dir:       config.Stop,
    ElevState: config.Idle,
    Queue:     [config.NumFloors][config.NumBtns]bool{},
  }

  //go orderhandler.AddOrdersToQueue(newOrder, &elevator)
  initState(&elevator)

  //Stop elevator in the first floor that the elevators arrive in
  for {
    select{
    case floor := <- ch.Drv_floors:
      elevator.Floor = floor
      elevio.SetMotorDirection(elevio.MD_Stop)
      elevio.SetFloorIndicator(floor)
      break
    }
    break
  }

  config.ElevatorMap[id] = elevator
  reciever <- config.ElevatorMap

  for {
    switch elevator.ElevState {
    case config.Idle:
        elevator.Dir = orderhandler.FindDirection(config.ElevatorMap[id])
        if elevator.Dir != config.Stop{
          if elevator.Dir == config.MovingUp{
            elevio.SetMotorDirection(elevio.MD_Up)
          }
          if elevator.Dir == config.MovingDown{
            elevio.SetMotorDirection(elevio.MD_Down)
          }
          elevator.ElevState = config.Moving
        }
        if orderhandler.CheckOrderSameFLoor(config.ElevatorMap[id], id){
          elevator.ElevState = config.ArrivedAtFloor
        }

    case config.Moving:
      select{
      case floor := <- ch.Drv_floors:
        elevio.SetFloorIndicator(floor)
        elevator.Floor = floor
        if orderhandler.CheckIfArrived(floor, config.ElevatorMap[id], id){
          elevator.ElevState = config.ArrivedAtFloor
        }
      }
    case config.ArrivedAtFloor:
      go timer.SetTimer(ch.Open_door, 3)
      reachedFloor(ch.Open_door, id)
    }
  }
}
