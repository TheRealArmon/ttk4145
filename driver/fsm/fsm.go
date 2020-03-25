package fsm

import "../elevio"
import "../config"
import "../timer"
import "../orderhandler"
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

func reachedFloor(sender <-chan bool, elevatorStatus *config.ElevatorState) {
  elevio.SetMotorDirection(elevio.MD_Stop)
  elevatorStatus.ElevState = config.Idle
  elevatorStatus.Dir = orderhandler.FindDirection(elevatorStatus)
  elevio.SetDoorOpenLamp(true)
  for {
    select{
    case <- sender:
      elevio.SetDoorOpenLamp(false)
      return
    }
  }
}

/*func updateElevatorList(elevator config.ElevatorState, id int, newState chan<- [config.NumElevators]config.ElevatorState,
  elevatorList *[config.NumElevators]config.ElevatorState){
  elevatorList[id] = elevator
  newState <- *elevatorList
  }*/


func ElevStateMachine(ch config.FSMChannels, id int, newState chan<- [config.NumElevators]config.ElevatorState,
  elevatorList *[config.NumElevators]config.ElevatorState) {


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
  //updateElevatorList(elevator, id, newState, elevatorList)

  changedState := true

  for {
    switch elevatorList[id].ElevState {
    case config.Idle:
        elevatorList[id].Dir = orderhandler.FindDirection(&elevatorList[id])
        if elevatorList[id].Dir != config.Stop{
          if elevatorList[id].Dir == config.MovingUp{
            elevio.SetMotorDirection(elevio.MD_Up)
            changedState = true
          }
          if elevatorList[id].Dir == config.MovingDown{
            elevio.SetMotorDirection(elevio.MD_Down)
            changedState = true
          }
          elevatorList[id].ElevState = config.Moving
        }
        if orderhandler.CheckOrderSameFLoor(&elevatorList[id], id){
          elevatorList[id].ElevState = config.ArrivedAtFloor
          changedState = true
        }
        if changedState{
          go func(){newState <- *elevatorList}()
          changedState = false
        }
        

    case config.Moving:
      select{
      case floor := <- ch.Drv_floors:
        elevio.SetFloorIndicator(floor)
        elevatorList[id].Floor = floor
        if orderhandler.CheckIfArrived(floor, &elevatorList[id], id){
          elevatorList[id].ElevState = config.ArrivedAtFloor
        }
        //newState <- *elevatorList
        //updateElevatorList(elevator, id, newState, elevatorList)
      }

    case config.ArrivedAtFloor:
      go timer.SetTimer(ch.Open_door, 3)
      reachedFloor(ch.Open_door, &elevatorList[id])
      newState <- *elevatorList
      //updateElevatorList(elevator, id, newState, elevatorList)
    }
  }
}


