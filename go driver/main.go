package main

import "./elevio"
import "./fsm"
import "./config"
import "./orderhandler"
import "./networkmod"
import "./networkmod/network/localip"
import "flag"
import "os"
import "fmt"


func main(){
    elevio.Init("localhost:15657", config.NumFloors)
    
    var id string
    flag.StringVar(&id, "id", "", "id of this peer")
    flag.Parse()

    // ... or alternatively, we can use the local IP address.
    // (But since we can run multiple programs on the same PC, we also append the
    //  process ID)
    if id == "" {
        localIP, err := localip.LocalIP()
        if err != nil {
            fmt.Println(err)
            localIP = "DISCONNECTED"
        }

        id = fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())

    }
    
    fsmChannels := config.FSMChannels{
      NewOrderToHandle: make(chan config.ElevatorOrder),
      Drv_buttons: make(chan elevio.ButtonEvent),
      Drv_floors: make(chan int),
      Drv_stop: make(chan bool),
      Open_door: make(chan bool),
    }
    
    go elevio.PollButtons(fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(fsmChannels.Drv_floors)
    go elevio.PollStopButton(fsmChannels.Drv_stop)
    go orderhandler.CheckNewOrder(fsmChannels.NewOrderToHandle, fsmChannels.Drv_buttons, 1)
    
    go networkmod.SendData(id)
    go networkmod.RecieveData(id)
    fsm.ElevStateMachine(fsmChannels, id)

}
