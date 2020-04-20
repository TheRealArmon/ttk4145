package fsm

import (
  "../elevio"
  cf "../config"
  "strconv"
  "time"
  "fmt"
)

func ElevStateMachine(driverCh cf.DriverChannels, id int, orderCh cf.OrderChannels,
    elevatorList *[cf.NumElevators]cf.ElevatorState, timerCh cf.TimerChannels) {

  idAsString := strconv.Itoa(id)
  idIndex := id - 1 //Each elevator's place in the elevator list corresponds to their id - 1
  
  //initilize the elevator and update state before sending state to peers on the network
  initState(&elevatorList[idIndex], driverCh.DrvFloors, id)
  go func(){orderCh.SendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()

  ticker := time.NewTicker(5 * time.Second)
  for {
    switch elevatorList[idIndex].State {
    
    case cf.Idle:
        elevatorList[idIndex].Dir = findDirection(&elevatorList[idIndex])
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
        if checkOrderSameFLoor(&elevatorList[idIndex]){
          elevatorList[idIndex].State = cf.ArrivedAtFloor
        }
        if (elevatorList[idIndex].State != cf.Idle){
          go func(){orderCh.SendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()
        }

    case cf.Moving:
      select{
      case floor := <- driverCh.DrvFloors:
        elevio.SetFloorIndicator(floor)
        elevatorList[idIndex].Floor = floor
        ticker.Stop()                          //Resets the ticker each time it moves past a floor
        ticker = time.NewTicker(5 * time.Second)
        if checkIfArrived(floor, &elevatorList[idIndex]){
          elevatorList[idIndex].State = cf.ArrivedAtFloor
        }
        go func(){orderCh.SendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()
      
      //Elevator has taken too long to get to a new floor which means that it has motor power loss 
      case <- ticker.C:
        ticker.Stop()
        elevatorList[idIndex].State = cf.SystemFailure
        orderCh.LostConnection <- elevatorList[idIndex]
        go func(){orderCh.SendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()
      }

    case cf.ArrivedAtFloor:
      //Send order letting the peers know that the order has been executed
      button := findOrderButton(elevatorList[idIndex].Floor, &elevatorList[idIndex])
      go func(){orderCh.SendOrder <- cf.ElevatorOrder{button, elevatorList[idIndex].Floor, id, true}}()
      
      ticker.Stop()
      clearOrderQueue(elevatorList[idIndex].Floor, &elevatorList[idIndex])
      go func(){time.Sleep(3 * time.Second); timerCh.Open_door <- true}()
      reachedFloor(timerCh.Open_door, &elevatorList[idIndex])
      if elevatorList[idIndex].State == cf.Moving{
        ticker =  time.NewTicker(5 * time.Second)
      }
      go func(){orderCh.SendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}}()

     //When the elevator starts moving again it stops at the first floor it arrives at
    case cf.SystemFailure:
      fmt.Println("System failure")
      select{
      case floor := <- driverCh.DrvFloors:
        elevatorList[idIndex].Floor = floor
        go func(){time.Sleep(3 * time.Second); timerCh.Open_door <- true}()
        reachedFloor(timerCh.Open_door, &elevatorList[idIndex])
      }
    }
  }
}
