package fsm

import "../elevio"
import "../config"
import "../timer"
import "../orderhandler"
import "strconv"
//import "sync"
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

func reachedFloor(door_timer <-chan bool, elevatorStatus *config.ElevatorState) {
  elevatorStatus.Dir = config.Stop
  orderhandler.SwitchOffButtonLight(elevatorStatus.Floor)
  elevio.SetMotorDirection(elevio.MD_Stop)
  elevio.SetDoorOpenLamp(true)
  for {
    select{
    case <- door_timer:
      elevio.SetDoorOpenLamp(false)
      elevatorStatus.Dir = orderhandler.FindDirection(elevatorStatus)
      if (elevatorStatus.Dir == config.Stop){
        elevatorStatus.ElevState = config.Idle
      }else{elevatorStatus.ElevState = config.Moving}
      setMotorDirection(elevatorStatus.Dir)
      return
    }
  }
}


func setMotorDirection(dir config.Directions){
  if (dir == config.MovingUp){
    elevio.SetMotorDirection(elevio.MD_Up)
  }
  if (dir == config.MovingDown){
    elevio.SetMotorDirection(elevio.MD_Down)
  }
}


func ElevStateMachine(ch config.FSMChannels, id int, sendOrder chan<- config.ElevatorOrder, newState chan<- map[string][config.NumElevators]config.ElevatorState,
  elevatorList *[config.NumElevators]config.ElevatorState, timerCh config.TimerChannels) {

  idAsString := strconv.Itoa(id)

  elevator := config.ElevatorState{
    Id:        id,
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

  elevatorList[id] = elevator

  for {
    switch elevatorList[id].ElevState {
    case config.Idle:
        elevatorList[id].Dir = orderhandler.FindDirection(&elevatorList[id])
        if elevatorList[id].Dir != config.Stop{
          if elevatorList[id].Dir == config.MovingUp{
            elevio.SetMotorDirection(elevio.MD_Up)
          }
          if elevatorList[id].Dir == config.MovingDown{
            elevio.SetMotorDirection(elevio.MD_Down)
          }
          elevatorList[id].ElevState = config.Moving
        }
        if orderhandler.CheckOrderSameFLoor(&elevatorList[id], id){
          elevatorList[id].ElevState = config.ArrivedAtFloor
        }
        if (elevatorList[id].ElevState != config.Idle){
          go func(){newState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}}()
        }

    case config.Moving:
      select{
      case floor := <- ch.Drv_floors:
        elevio.SetFloorIndicator(floor)
        elevatorList[id].Floor = floor
        if orderhandler.CheckIfArrived(floor, &elevatorList[id], id){
          elevatorList[id].ElevState = config.ArrivedAtFloor
        }
        go func(){newState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}}()
      }

    case config.ArrivedAtFloor:
      go func(){sendOrder <- config.ElevatorOrder{elevio.BT_HallUp, elevatorList[id].Floor, id, true}}()
      go func(){sendOrder <- config.ElevatorOrder{elevio.BT_HallDown, elevatorList[id].Floor, id, true}}()
      go timer.SetTimer(timerCh, config.Door)
      reachedFloor(timerCh.Open_door, &elevatorList[id])
      go func(){newState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}}()
    }
  }
}
