package orderhandler

import (
	"time"

	"./../elevio"
)

//import "fmt"
//import "net"

type OrderType struct {
	Floor int
	//Button ButtonType
}

var _numFloors int = 4

/*func handle_order(event elevio.ButtonEvent, floor int){
  _mtx.Lock()
	defer _mtx.Unlock()
  if event.Floor < floor {
    elevio.SetMotorDirection(elevio.MD_Down)
  }
  if event.Floor > floor{
    elevio.SetMotorDirection(elevio.MD_Up)
  }
}

func stop_elev(event elevio.ButtonEvent, floor int){
  _mtx.Lock()
	defer _mtx.Unlock()
  if (event.Floor == floor){
    fmt.Printf("nå stopper heisen ikke")
    elevio.SetMotorDirection(elevio.MD_Stop)
  }
}*/

func CheckFloorOrder(reciever chan<- OrderType) {
	prev := make([][3]bool, _numFloors)
	for {
		time.Sleep(20 * time.Millisecond)
		for f := 0; f < _numFloors; f++ {
			for b := elevio.ButtonType(0); b < 3; b++ {
				v := elevio.GetButton(b, f)
				if v != prev[f][b] && v != false {
					reciever <- OrderType{f}
				}
				prev[f][b] = v
			}
		}
	}
}
