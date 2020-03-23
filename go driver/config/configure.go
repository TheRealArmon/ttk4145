package config

import "../elevio"
import "../networkmod/network/peers"
import "os"
import "sync"

var LOCAL_ID = os.Args[1]

const (
  NumFloors int = 4
  StartFloor    = 2
  NumElevators  = 3
  NumBtns       = 3
)



var ElevatorMap = make(map[string]ElevatorState)
var ActiveElevatorMap = make(map[string]bool)
var mutex = &sync.Mutex{}


 



type State int

const(
  Idle State      = 0
  Moving          = 1
  ArrivedAtFloor  = 2
)

type Directions int

const(
  MovingDown Directions = -1
  Stop                  = 0
  MovingUp              = 1
)

type ElevatorState struct{
  Floor                       int
  Dir                         Directions
  ElevState                   State
  Queue [NumFloors][NumBtns]  bool
}

type FSMChannels struct {
  Drv_buttons       chan elevio.ButtonEvent
  Drv_floors        chan int
  Drv_stop          chan bool
  Open_door         chan bool
}



type NetworkChannels struct {
    PeerTxEnable          chan bool
    PeerUpdateCh          chan peers.PeerUpdate
    TransmittOrderCh      chan ElevatorOrder
    TransmittStateCh      chan map[string]ElevatorState
    RecieveOrderCh        chan ElevatorOrder
    RecieveStateCh        chan map[string]ElevatorState
}

type ElevatorOrder struct{
  Button              elevio.ButtonType
  Floor               int
  ExecutingElevator   string
  OrderStatus         bool
}



