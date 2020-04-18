package fsm

import "../elevio"
import "../config"
import "../timer"
import "../orderhandler"
import "strconv"
//import "fmt"
import "time"

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
      if orderhandler.CheckOrderSameFLoor(elevatorStatus){
        elevatorStatus.ElevState = config.ArrivedAtFloor
        return
      }
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


func ElevStateMachine(ch config.FSMChannels, id int, sendOrder chan<- config.ElevatorOrder, sendState chan<- map[string][config.NumElevators]config.ElevatorState,
  elevatorList *[config.NumElevators]config.ElevatorState, timerCh config.TimerChannels, lostConnection chan<- config.ElevatorState) {

  idAsString := strconv.Itoa(id)
  idIndex := id - 1

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

  elevatorList[idIndex] = elevator
  sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
  ticker := time.NewTicker(5000 * time.Millisecond)
  for {
    switch elevatorList[idIndex].ElevState {
    case config.Idle:
        elevatorList[idIndex].Dir = orderhandler.FindDirection(&elevatorList[idIndex])
        if elevatorList[idIndex].Dir != config.Stop{
          if elevatorList[idIndex].Dir == config.MovingUp{
            elevio.SetMotorDirection(elevio.MD_Up)
          }
          if elevatorList[idIndex].Dir == config.MovingDown{
            elevio.SetMotorDirection(elevio.MD_Down)
          }
          elevatorList[idIndex].ElevState = config.Moving
          ticker = time.NewTicker(5000 * time.Millisecond)
        }
        if orderhandler.CheckOrderSameFLoor(&elevatorList[idIndex]){
          elevatorList[idIndex].ElevState = config.ArrivedAtFloor
        }
        if (elevatorList[idIndex].ElevState != config.Idle){
          sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
        }

    case config.Moving:
      select{
      case floor := <- ch.Drv_floors:
        ticker.Stop()
        ticker = time.NewTicker(5000 * time.Millisecond)
        elevio.SetFloorIndicator(floor)
        elevatorList[idIndex].Floor = floor
        if orderhandler.CheckIfArrived(floor, &elevatorList[idIndex]){
          elevatorList[idIndex].ElevState = config.ArrivedAtFloor
        }
        sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
      case <- ticker.C:
        ticker.Stop()
        elevatorList[idIndex].ElevState = config.SystemFailure
        lostConnection <- elevatorList[idIndex]
        sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
      }

    case config.ArrivedAtFloor:
      ticker.Stop()
      sendOrder <- config.ElevatorOrder{elevio.BT_HallUp, elevatorList[idIndex].Floor, id, true}
      sendOrder <- config.ElevatorOrder{elevio.BT_HallDown, elevatorList[idIndex].Floor, id, true}
      sendOrder <- config.ElevatorOrder{elevio.BT_Cab, elevatorList[idIndex].Floor, id, true}
      go timer.SetTimer(timerCh, config.Door)
      reachedFloor(timerCh.Open_door, &elevatorList[idIndex])
      if elevatorList[idIndex].ElevState == config.Moving{
        ticker =  time.NewTicker(5000 * time.Millisecond)
      }
      sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}

    case config.SystemFailure:
      
    }
  }
}
