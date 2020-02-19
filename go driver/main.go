package main

import "./elevio"
import "fmt"
//import "./orderhandler"




func main(){

    numFloors := 4
    elevio.Init("localhost:15657", numFloors)

    var current_direction elevio.MotorDirection = elevio.MD_Stop
    var current_floor = 0
    var order_floor = 0
    var order_type elevio.ButtonType = elevio.BT_Cab


    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)
  //  drv_order   := make(chan OrderType)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)
    //go orderhandler.CheckFloorOrder(drv_order)


    for {
        select {
        case a := <- drv_buttons:
            order_floor = a.Floor
            order_type = a.Button

            if current_direction == elevio.MD_Stop{

              if order_floor == current_floor{
                current_direction = elevio.MD_Stop
              }
              if order_floor < current_floor {
                elevio.SetMotorDirection(elevio.MD_Down)
                current_direction = elevio.MD_Down
                elevio.SetButtonLamp(order_type, order_floor, true)
              }
              if order_floor > current_floor{
                elevio.SetMotorDirection(elevio.MD_Up)
                current_direction = elevio.MD_Up
                elevio.SetButtonLamp(order_type, order_floor, true)
               }
            }

        case a := <- drv_floors:
            current_floor = a
            elevio.SetFloorIndicator(current_floor)
            if a == order_floor {
              elevio.SetMotorDirection(elevio.MD_Stop)
              current_direction = elevio.MD_Stop
              elevio.SetButtonLamp(order_type,order_floor,false)
            }

            if a == numFloors-1 && order_floor != numFloors-1 {
                current_direction = elevio.MD_Down
            } else if a == 0 && order_floor != 0{
                current_direction = elevio.MD_Up
            }
            elevio.SetMotorDirection(current_direction)


        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(current_direction)
            }

        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < numFloors; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        //case a := <- drv_order:



        }
    }
}
