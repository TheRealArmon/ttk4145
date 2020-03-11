package main

import "./elevio"
import "./fsm"
import "./config"
import "./orderhandler"
import "./networkmod"
import "./networkmod/network/localip"
import "./networkmod/network/peers"
import "./networkmod/network/bcaste"
import "flag"
import "os"
import "fmt"
e

func main(){
    elevio.Init("localhost:12345", config.NumFloors)
    
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
    
    networkChannels := config.NetworkChannels{
      PeerTxEnable: make(chan bool),
      PeerUpdateCh: make(chan peers.PeerUpdate),
      TransmitterCh: make(chan config.ElevatorState),
      RecieveCh: make(chan config.ElevatorState),
    }
    
    go peers.Transmitter(12346, id, networkChannels.PeerTxEnable)
    go peers.Receiver(12346, networkChannels.PeerUpdateCh)
    go bcast.Transmitter(12347, networkChannels.TransmitterCh)
    go bcast.Receiver(12347, networkChannels.RecieveCh)
    
    go elevio.PollButtons(fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(fsmChannels.Drv_floors)
    go elevio.PollStopButton(fsmChannels.Drv_stop)
    go orderhandler.CheckNewOrder(fsmChannels.NewOrderToHandle, fsmChannels.Drv_buttons, 1)
    
    go networkmod.SendData(id, networkChannels)
    //go networkmod.RecieveData(id)
    fsm.ElevStateMachine(fsmChannels, id)

}
