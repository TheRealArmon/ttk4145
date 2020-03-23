package main

import "./elevio"
import "./fsm"
import "./config"
import "./orderhandler"
import "./networkmod"
import "./networkmod/network/localip"
import "./networkmod/network/peers"
import "./networkmod/network/bcast"
import "os"
import "fmt"
import "flag"


func main(){
    var Local_Id string
    flag.StringVar(&Local_Id, "id", "", "Id of this peer")
    flag.Parse()


    elevio.Init("localhost:"+Local_Id, config.NumFloors)

    // ... or alternatively, we can use the local IP address.
    // (But since we can run multiple programs on the same PC, we also append the
    //  process ID)

    localIP, err := localip.LocalIP()
    if err != nil {
        fmt.Println(err)
        localIP = "DISCONNECTED"
    }
    id := fmt.Sprintf("peer-%s-%d", localIP, os.Getpid())


    fsmChannels := config.FSMChannels{
      Drv_buttons: make(chan elevio.ButtonEvent),
      Drv_floors: make(chan int),
      Drv_stop: make(chan bool),
      Open_door: make(chan bool),
    }

    newOrder := make(chan config.ElevatorOrder)
    newState := make(chan map[string]config.ElevatorState)


    networkChannels := config.NetworkChannels{
      PeerTxEnable: make(chan bool),
      PeerUpdateCh: make(chan peers.PeerUpdate),
      TransmittOrderCh: make(chan config.ElevatorOrder),
      TransmittStateCh: make(chan map[string]config.ElevatorState),
      RecieveOrderCh: make(chan config.ElevatorOrder),
      RecieveStateCh: make(chan map[string]config.ElevatorState),
    }


    go peers.Transmitter(12346, id, networkChannels.PeerTxEnable)
    go peers.Receiver(12346, networkChannels.PeerUpdateCh)
    go bcast.Transmitter(12347, networkChannels.TransmittOrderCh)
    go bcast.Receiver(12347, networkChannels.RecieveOrderCh)
    go bcast.Transmitter(12348, networkChannels.TransmittStateCh)
    go bcast.Receiver(12348, networkChannels.RecieveStateCh)

    go elevio.PollButtons(fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(fsmChannels.Drv_floors)
    go elevio.PollStopButton(fsmChannels.Drv_stop)
    //go orderhandler.CheckNewOrder(newOrder, fsmChannels.Drv_buttons, id)


    go networkmod.SendData(id, networkChannels, newOrder, newState)
    go networkmod.RecieveData(id, networkChannels)

    go orderhandler.OrderHandler(fsmChannels.Drv_buttons, newOrder, id)
    fsm.ElevStateMachine(fsmChannels, newOrder, id, newState)

}
