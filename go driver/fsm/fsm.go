package fsm

import "../elevio"
import "../config"
import "../timer"
import "../orderhandler"
import "sync"
//import "fmt"
import "time"
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
  elevatorStatus.Dir = config.Stop
  elevio.SetDoorOpenLamp(true)
  for {
    select{
    case <- sender:
      elevio.SetDoorOpenLamp(false)
      return
    }
  }
}

func updateElevatorMap(elevator config.ElevatorState, mutex *sync.RWMutex, id string, newState chan<- map[string]config.ElevatorState){
  mutex.Lock()
  config.ElevatorMap[id] = elevator
  var temp = config.ElevatorMap
  newState <- temp
  mutex.Unlock()
  }


func ElevStateMachine(ch config.FSMChannels, id string, newState chan<- map[string]config.ElevatorState,
  mutex *sync.RWMutex) {


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


  go updateElevatorMap(elevator, mutex, id, newState)

  for {
    switch elevator.ElevState {
    case config.Idle:
        elevator.Dir = orderhandler.FindDirection(&elevator)
        if elevator.Dir != config.Stop{
          if elevator.Dir == config.MovingUp{
            elevio.SetMotorDirection(elevio.MD_Up)
          }
          if elevator.Dir == config.MovingDown{
            elevio.SetMotorDirection(elevio.MD_Down)
          }
          elevator.ElevState = config.Moving
        }
        if orderhandler.CheckOrderSameFLoor(&elevator, id){
          elevator.ElevState = config.ArrivedAtFloor
        }
        go updateElevatorMap(elevator, mutex, id, newState)
        time.Sleep(10 * time.Millisecond)

    case config.Moving:
      select{
      case floor := <- ch.Drv_floors:
        elevio.SetFloorIndicator(floor)
        elevator.Floor = floor
        if orderhandler.CheckIfArrived(floor, &elevator, id){
          elevator.ElevState = config.ArrivedAtFloor
        }
        go updateElevatorMap(elevator, mutex, id, newState)
      }

    case config.ArrivedAtFloor:
      go timer.SetTimer(ch.Open_door, 3)
      reachedFloor(ch.Open_door, &elevator)
      go updateElevatorMap(elevator, mutex, id, newState)
    }
  }
}


