package config

import "../elevio"
import "../networkmod/network/peers"
import "os"

var LOCAL_ID = os.Args[1]

const (
  NumFloors int = 4
  NumElevators  = 3
  NumBtns       = 3
)

type States int

const(
  Idle States      = 0
  Moving          = 1
  ArrivedAtFloor  = 2
  SystemFailure   = 3
)

type Directions int

const(
  MovingDown Directions = -1
  Stop                  = 0
  MovingUp              = 1
)


type TimerCase int

const(
  Door TimerCase = 0
)

type TimerChannels struct {
  Open_door     chan bool
}

type ElevatorState struct{
  Id                          int
  Floor                       int
  Dir                         Directions
  State                       States
  Queue [NumFloors][NumBtns]  bool
}

type DriverChannels struct {
  DrvButtons       chan elevio.ButtonEvent
  DrvFloors        chan int
  DrvStop          chan bool
}

type NetworkChannels struct {
    PeerTxEnable          chan bool
    PeerUpdateCh          chan peers.PeerUpdate
    TransmittOrderCh      chan ElevatorOrder
    TransmittStateCh      chan map[string][NumElevators]ElevatorState
    RecieveOrderCh        chan ElevatorOrder
    RecieveStateCh        chan map[string][NumElevators]ElevatorState
}

type OrderChannels struct {
  LostConnection  chan ElevatorState
  SendState       chan map[string][NumElevators]ElevatorState
  SendOrder       chan ElevatorOrder
}

type ElevatorOrder struct{
  Button              elevio.ButtonType
  Floor               int
  ExecutingElevator   int
  OrderStatus         bool
}



