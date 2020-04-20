package main

import (
  "./elevio"
  "./fsm"
  "./config"
  "./orderhandler"
  "./networkmod"
  "./networkmod/network/peers"
  "./networkmod/network/bcast"
  "flag"
  "strconv"
)

func main(){
  
  var (
    localHostId        string
    id                 int
    ElevatorList       [config.NumElevators]config.ElevatorState
    ActiveElevatorList [config.NumElevators]bool
  )

    flag.StringVar(&localHostId, "hostId", "", "hostId of this peer")
    flag.IntVar(&id, "id", 1234, "id of this peer")
    flag.Parse()
    elevio.Init("localhost:"+localHostId, config.NumFloors)

    idAsString := strconv.Itoa(id)


    driverChannels := config.DriverChannels{
      DrvButtons: make(chan elevio.ButtonEvent),
      DrvFloors: make(chan int),
      DrvStop: make(chan bool),
    }

    networkChannels := config.NetworkChannels{
      PeerTxEnable: make(chan bool),
      PeerUpdateCh: make(chan peers.PeerUpdate),
      TransmittOrderCh: make(chan config.ElevatorOrder),
      TransmittStateCh: make(chan map[string][config.NumElevators]config.ElevatorState),
      RecieveOrderCh: make(chan config.ElevatorOrder),
      RecieveStateCh: make(chan map[string][config.NumElevators]config.ElevatorState),
    }

    timerChannels := config.TimerChannels{
      Open_door: make(chan bool),
    }


    orderChannels := config.OrderChannels{
      LostConnection: make(chan config.ElevatorState),
      SendState: make(chan map[string][config.NumElevators]config.ElevatorState),
      SendOrder: make(chan config.ElevatorOrder),
    }

    go peers.Transmitter(22349, idAsString, networkChannels.PeerTxEnable)
    go peers.Receiver(22349, networkChannels.PeerUpdateCh)
    go bcast.Transmitter(22367, networkChannels.TransmittOrderCh)
    go bcast.Receiver(22367, networkChannels.RecieveOrderCh)
    go bcast.Transmitter(22378, networkChannels.TransmittStateCh)
    go bcast.Receiver(22378, networkChannels.RecieveStateCh)

    go elevio.PollButtons(driverChannels.DrvButtons)
    go elevio.PollFloorSensor(driverChannels.DrvFloors)
    go elevio.PollStopButton(driverChannels.DrvStop)

    go networkmod.RecieveData(id, networkChannels, orderChannels.LostConnection, &ElevatorList, &ActiveElevatorList)
    go networkmod.SendData(networkChannels, orderChannels) 

    go orderhandler.OrderHandler(driverChannels.DrvButtons, orderChannels.SendOrder, orderChannels.SendState,
      networkChannels.RecieveStateCh, networkChannels.RecieveOrderCh, orderChannels.LostConnection, id, &ElevatorList, &ActiveElevatorList)
    fsm.ElevStateMachine(driverChannels, id, orderChannels.SendOrder, orderChannels.SendState, &ElevatorList, timerChannels, orderChannels.LostConnection)

}