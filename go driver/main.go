package main

import "./elevio"
import "fmt"
import "sync"

var _mtx sync.Mutex

func handle_order(event elevio.ButtonEvent, floor int){
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
}

func main(){

    numFloors := 4
    elevio.Init("localhost:15657", numFloors)

    var d elevio.MotorDirection = elevio.MD_Up
//elevio.SetMotorDirection(d)

    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)


    for {
        select {
        case a := <- drv_buttons:
            //fmt.Printf("%+v\n", a)
            b := <- drv_floors
            //fmt.Printf("%+v", b)
            elevio.SetButtonLamp(a.Button, a.Floor, true)
            handle_order(a, b)

        case a := <- drv_floors:
            b := <- drv_buttons
            fmt.Printf("%+v\n",b.Floor)
            stop_elev(b, a)
            /*if a == numFloors-1 {
                d = elevio.MD_Down
            } else if a == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)*/


        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }

        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        }
    }
}
