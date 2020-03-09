package fsm

import "../elevio"
import "../config"
import "../timer"
import "../orderhandler"
//import "time"
import "fmt"

var _currentDirection elevio.MotorDirection
var _destination int
var _current_direction elevio.MotorDirection
var _current_floor int
var _order_type elevio.ButtonType
var _orderQueue [config.NumFloors][config.NumBtns] bool


func initState(elevator *config.ElevatorState) {
  for i := 0; i < config.NumFloors; i++{
    for j := elevio.BT_HallUp; j < config.NumBtns; j++{
      elevio.SetButtonLamp(j, i, false)
      elevator.Queue[i][j] = false
    }
  }
  elevio.SetMotorDirection(elevio.MD_Down)
}

func reachedFloor(sender <-chan bool) {
  elevio.SetMotorDirection(elevio.MD_Stop)
  elevio.SetButtonLamp(_order_type,_destination,false)
  _current_direction = elevio.MD_Stop
  elevio.SetDoorOpenLamp(true)
  for {
    select{
    case <- sender:
      elevio.SetDoorOpenLamp(false)
      return
    }
  }
}


//SKRIV OM ELLER IKKE??

// func checkReachedEdges() {
//   if _current_floor == config.NumFloors-1 && _destination != config.NumFloors-1 {
//       _current_direction = elevio.MD_Stop
//       elevio.SetMotorDirection(_current_direction)
//   } else if _current_floor == 0 && _destination != 0{
//       _current_direction = elevio.MD_Stop
//       elevio.SetMotorDirection(_current_direction)
//   }
// }



func ElevStateMachine(ch config.FSMChannels) {

  elevator := config.ElevatorState{
    ID: 1,
    Dir:       config.Stop,
    ElevState: config.Idle,
    Queue:     [config.NumFloors][config.NumBtns]bool{},
  }
  go orderhandler.AddOrdersToQueue(ch.NewOrderToHandle, &elevator)

  initState(&elevator)

  //Stop elevator in the first floor that the elevators arrive in
  for {
    select{
    case floor := <- ch.Drv_floors:
      elevator.Floor = floor
      elevio.SetMotorDirection(elevio.MD_Stop)
      go timer.SetTimer(ch.Close_door, 3)
      reachedFloor(ch.Close_door)
      break
    }
    break
  }

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

    case config.Moving:
      go CheckIfArrived(ch.Drv_floors, &elevator)
    case config.ArrivedAtFloor:
      PrintQueu(&elevator)
      elevio.SetMotorDirection(elevio.MD_Stop)
      elevator.ElevState = config.Idle
      elevator.Dir = config.Stop
    }
  }
}

//Sjekk om det er ordre i køen

func CheckIfArrived(sender <-chan int, elevator *config.ElevatorState){
  for{
    select{
    case floor := <- sender:
      for i := 0; i < config.NumFloors; i++{
        for j := elevio.BT_HallUp; j < config.NumBtns; j++{
          if i == floor{
            elevator.Queue[i][j] = false
            elevator.ElevState = config.ArrivedAtFloor
        }
      }
    }
    }
  }
}



// for {
//     select {
//     case newOrder := <- ch.NewOrderToHandle:
//         order_floor := newOrder.Floor
//         button_type := newOrder.Button
//         elevator.Queue[order_floor][button_type] = true
//         elevator.Dir = orderhandler.FindDirection(elevator)
//
//     case a := <- ch.Drv_floors:
//       switch elevator.ElevState {
//       case config.Init:
//         elevio.SetMotorDirection(elevio.MD_Stop)
//         elevator.ElevState = config.Idle
//       }
//         _current_floor = a
//         elevio.SetFloorIndicator(_current_floor)
//         for i := 0; i<3; i++{
//           if (_orderQueue[_current_floor][i] == true){
//             _orderQueue[_current_floor][i] = false
//             go timer.SetTimer(ch.Close_door, 3)
//             reachedFloor(ch.Close_door)
//           }
//         }
//         checkReachedEdges()
//     }
// }
