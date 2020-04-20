package fsm

import (
  "../elevio"
  cf "../config"
  "../timer"
  oh "../orderhandler"
  "strconv"
  "time"
  "fmt"
)

func ElevStateMachine(ch cf.DriverChannels, id int, sendOrder chan<- cf.ElevatorOrder, sendState chan<- map[string][cf.NumElevators]cf.ElevatorState,
  elevatorList *[cf.NumElevators]cf.ElevatorState, timerCh cf.TimerChannels, lostConnection chan<- cf.ElevatorState) {

  idAsString := strconv.Itoa(id)
  idIndex := id - 1 //Each elevator's place in the elevator list corresponds to their id - 1
  
  //initilize the elevator and update state before sending state to peers on the network
  initState(&elevatorList[idIndex], ch.DrvFloors, id)
  go func(){sendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()

  ticker := time.NewTicker(5 * time.Second)
  for {
    switch elevatorList[idIndex].State {
    
    case cf.Idle:
        elevatorList[idIndex].Dir = oh.FindDirection(&elevatorList[idIndex])
        if elevatorList[idIndex].Dir != cf.Stop{
          if elevatorList[idIndex].Dir == cf.MovingUp{
            elevio.SetMotorDirection(elevio.MD_Up)
          }
          if elevatorList[idIndex].Dir == cf.MovingDown{
            elevio.SetMotorDirection(elevio.MD_Down)
          }
          elevatorList[idIndex].State = cf.Moving
          ticker = time.NewTicker(5 * time.Second)  //Start ticker as the elevator starts moving to detect power loss
        }
        if oh.CheckOrderSameFLoor(&elevatorList[idIndex]){
          elevatorList[idIndex].State = cf.ArrivedAtFloor
        }
        if (elevatorList[idIndex].State != cf.Idle){
          go func(){sendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()
        }

    case cf.Moving:
      select{
      case floor := <- ch.DrvFloors:
        elevio.SetFloorIndicator(floor)
        ticker.Stop()
        ticker = time.NewTicker(5 * time.Second)
        elevatorList[idIndex].Floor = floor
        if oh.CheckIfArrived(floor, &elevatorList[idIndex]){
          elevatorList[idIndex].State = cf.ArrivedAtFloor
        }
        go func(){sendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()
      
      case <- ticker.C:
        ticker.Stop()
        elevatorList[idIndex].State = cf.SystemFailure
        lostConnection <- elevatorList[idIndex]
        go func(){sendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()
      }

    case cf.ArrivedAtFloor:
      ticker.Stop()
      button := oh.FindOrderButton(elevatorList[idIndex].Floor, &elevatorList[idIndex])
      go func(){sendOrder <- cf.ElevatorOrder{button, elevatorList[idIndex].Floor, id, true}}()
      oh.ClearOrderQueue(elevatorList[idIndex].Floor, &elevatorList[idIndex])
      go timer.SetTimer(timerCh, cf.Door)
      reachedFloor(timerCh.Open_door, &elevatorList[idIndex])
      if elevatorList[idIndex].State == cf.Moving{
        ticker =  time.NewTicker(5000 * time.Millisecond)
      }
      go func(){sendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()

    case cf.SystemFailure:
      fmt.Println("System failure")
      select{
      case floor := <- ch.DrvFloors:
        elevatorList[idIndex].Floor = floor
        go timer.SetTimer(timerCh, cf.Door)
        reachedFloor(timerCh.Open_door, &elevatorList[idIndex])
      }
    }
  }
}
