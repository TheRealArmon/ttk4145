package main

import "./elevio"
import "./fsm"
import "./config"
import "./orderhandler"
import "./networkmod"
import "./networkmod/network/peers"
import "./networkmod/network/bcast"
import "flag"
import "strconv"

func main(){

    var Local_Host_Id string
    var id int
    flag.StringVar(&Local_Host_Id, "hostId", "", "hostId of this peer")
    flag.IntVar(&id, "id", 1234, "id of this peer")
    flag.Parse()

    elevio.Init("localhost:"+Local_Host_Id, config.NumFloors)

    idAsString := strconv.Itoa(id)

    fsmChannels := config.FSMChannels{
      Drv_buttons: make(chan elevio.ButtonEvent),
      Drv_floors: make(chan int),
      Drv_stop: make(chan bool),
    }

    timerChannels := config.TimerChannels{
      Open_door: make(chan bool),
    }


    var ElevatorList [config.NumElevators]config.ElevatorState
    var ActiveElevatorList [config.NumElevators]bool


    lostConnection := make(chan config.ElevatorState)
    sendState := make(chan map[string][config.NumElevators]config.ElevatorState)
    sendOrder := make(chan config.ElevatorOrder)

    networkChannels := config.NetworkChannels{
      PeerTxEnable: make(chan bool),
      PeerUpdateCh: make(chan peers.PeerUpdate),
      TransmittOrderCh: make(chan config.ElevatorOrder),
      TransmittStateCh: make(chan map[string][config.NumElevators]config.ElevatorState),
      RecieveOrderCh: make(chan config.ElevatorOrder),
      RecieveStateCh: make(chan map[string][config.NumElevators]config.ElevatorState),
    }

    go peers.Transmitter(22349, idAsString, networkChannels.PeerTxEnable)
    go peers.Receiver(22349, networkChannels.PeerUpdateCh)
    go bcast.Transmitter(22367, networkChannels.TransmittOrderCh)
    go bcast.Receiver(22367, networkChannels.RecieveOrderCh)
    go bcast.Transmitter(22378, networkChannels.TransmittStateCh)
    go bcast.Receiver(22378, networkChannels.RecieveStateCh)

    go elevio.PollButtons(fsmChannels.Drv_buttons)
    go elevio.PollFloorSensor(fsmChannels.Drv_floors)
    go elevio.PollStopButton(fsmChannels.Drv_stop)

    go networkmod.RecieveData(id, networkChannels, &ElevatorList, &ActiveElevatorList, lostConnection)
    go networkmod.SendData(networkChannels, sendOrder, sendState) 

    go orderhandler.OrderHandler(fsmChannels.Drv_buttons, sendOrder, sendState,
      networkChannels.RecieveStateCh, networkChannels.RecieveOrderCh, lostConnection, id, &ElevatorList, &ActiveElevatorList)
    fsm.ElevStateMachine(fsmChannels, id, sendOrder, sendState, &ElevatorList, timerChannels, lostConnection)

}