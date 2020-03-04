package fsm

import "../elevio"
import "../config"
import "../timer"
import "fmt"

var _currentDirection elevio.MotorDirection
var _destination int
var _current_direction elevio.MotorDirection
var _current_floor int
var _order_type elevio.ButtonType
var _orderQueue [config.NumFloors][config.NumBtns] bool


func initState() {
  _destination = config.StartFloor
  _currentDirection = elevio.MD_Down
  elevio.SetMotorDirection(elevio.MD_Down)
  for i := 0; i < config.NumFloors; i++{
    for j := 0; j < config.NumBtns; j++{
      _orderQueue[i][j] = false
    }
  }
  _orderQueue[0][0] = true
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

func checkReachedEdges() {
  if _current_floor == config.NumFloors-1 && _destination != config.NumFloors-1 {
      _current_direction = elevio.MD_Stop
      elevio.SetMotorDirection(_current_direction)
  } else if _current_floor == 0 && _destination != 0{
      _current_direction = elevio.MD_Stop
      elevio.SetMotorDirection(_current_direction)
  }
}

func ElevStateMachine(ch config.FSMChannels) {


  initState()

  elevator := config.ElevatorState{
    ID: 1,
    ElevState: config.Idle,
    Floor:     config.StartFloor,
    Dir:       config.Stop,
    Queue:     _orderQueue,
  }

  for {
      select {
      case newOrder := <- ch.NewOrderToHandle:
          order_floor := newOrder.Floor
          button_type := newOrder.Button
          elevator.Queue[order_floor][button_type] = true
          elevator.Dir = FindDirection(elevator)
          fmt.Println("+%v", elevator.Dir)
          fmt.Println("+%v", elevator.Floor)

      case a := <- ch.Drv_floors:
          _current_floor = a
          elevio.SetFloorIndicator(_current_floor)
          for i := 0; i<3; i++{
            if (_orderQueue[_current_floor][i] == true){
              _orderQueue[_current_floor][i] = false
              go timer.SetDoorTimer(ch.Close_door)
              reachedFloor(ch.Close_door)
            }
          }
          checkReachedEdges()
      }
  }
}

func handleOrders(){

}
