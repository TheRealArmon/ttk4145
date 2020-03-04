package main

import "./elevio"
import "./fsm"
import "./config"
//import "./orderhandler"




func main(){
    elevio.Init("localhost:15657", config.NumFloors)
    fsmChannels := config.FSMChannels{
      Drv_buttons: make(chan elevio.ButtonEvent),
      Drv_floors: make(chan int),
      Drv_stop: make(chan bool),
      Close_door: make(chan bool),
    }

    go elevio.PollButtons(fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(fsmChannels.Drv_floors)
    go elevio.PollStopButton(fsmChannels.Drv_stop)

    fsm.ElevStateMachine(fsmChannels)

}
