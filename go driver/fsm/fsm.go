package fsm

//import "fmt"
import "../elevio"
import "../config"
import "../timer"

var _currentDirection elevio.MotorDirection
var _destination int
var _current_direction elevio.MotorDirection
var _current_floor int
var _order_type elevio.ButtonType

//lag funksjon for å åpne dør når en ankommer en etasje/trykker på knapp i samme etasje
//lag et internt køsystem
//lag initfunksjon som sender heisen til første etasje

func setDirection(order_type elevio.ButtonType) {
    // if _destination == _current_floor{
    //   _current_direction = elevio.MD_Stop
    //   //kjør funsksjon som åpner døren i 3 sek
    // }
    if _destination < _current_floor {
      elevio.SetMotorDirection(elevio.MD_Down)
      _current_direction = elevio.MD_Down
      elevio.SetButtonLamp(order_type, _destination, true)
    }
    if _destination > _current_floor{
      elevio.SetMotorDirection(elevio.MD_Up)
      _current_direction = elevio.MD_Up
      elevio.SetButtonLamp(order_type, _destination, true)
     }
}

func initState() {
  _destination = config.StartFloor
  _currentDirection = elevio.MD_Down
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

func checkReachedEdges() {
  if _current_floor == config.NumFloors-1 && _destination != config.NumFloors-1 {
      _current_direction = elevio.MD_Stop
      elevio.SetMotorDirection(_current_direction)
  } else if _current_floor == 0 && _destination != 0{
      _current_direction = elevio.MD_Stop
      elevio.SetMotorDirection(_current_direction)
  }
}

func ElevStateMachine() {
  drv_buttons := make(chan elevio.ButtonEvent)
  drv_floors  := make(chan int)
  drv_obstr   := make(chan bool)
  drv_stop    := make(chan bool)
  close_door  := make(chan bool)

  go elevio.PollButtons(drv_buttons)
  go elevio.PollFloorSensor(drv_floors)
  go elevio.PollObstructionSwitch(drv_obstr)
  go elevio.PollStopButton(drv_stop)

  initState()

  for {
      select {
      case a := <- drv_buttons:
          _destination = a.Floor
          _order_type = a.Button
          if _destination == _current_floor{
            go timer.SetDoorTimer(close_door)
            reachedFloor(close_door)
          } else {
            setDirection(_order_type)
          }

      case a := <- drv_floors:
          _current_floor = a
          elevio.SetFloorIndicator(_current_floor)
          if a == _destination || _current_floor == _destination{
            go timer.SetDoorTimer(close_door)
            reachedFloor(close_door)
          }
          checkReachedEdges()
      // case a := <- drv_stop:
      //     fmt.Printf("%+v\n", a)
      //     for f := 0; f < config.NumFloors; f++ {
      //         for b := elevio.ButtonType(0); b < 3; b++ {
      //             elevio.SetButtonLamp(b, f, false)
      //         }
      //     }
      }
  }
}
