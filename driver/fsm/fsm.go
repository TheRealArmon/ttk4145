package fsm

<<<<<<< HEAD
import (
	"strconv"

	"../config"
	"../elevio"
	"../orderhandler"
	"../timer"
)

//import "sync"
//import "fmt"
//import "time"
//import "reflect"
=======
import "../elevio"
import "../config"
import "../timer"
import "../orderhandler"
import "strconv"
>>>>>>> origin/development

func initState(elevator *config.ElevatorState) {
	elevio.SetDoorOpenLamp(false)
	for i := 0; i < config.NumFloors; i++ {
		for j := elevio.BT_HallUp; j < config.NumBtns; j++ {
			elevio.SetButtonLamp(j, i, false)
			elevator.Queue[i][j] = false
		}
	}
	elevio.SetMotorDirection(elevio.MD_Down)
}

func reachedFloor(door_timer <-chan bool, elevatorStatus *config.ElevatorState) {
	elevatorStatus.Dir = config.Stop
	orderhandler.SwitchOffButtonLight(elevatorStatus.Floor)
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	for {
		select {
		case <-door_timer:
			elevio.SetDoorOpenLamp(false)
			elevatorStatus.Dir = orderhandler.FindDirection(elevatorStatus)
			if elevatorStatus.Dir == config.Stop {
				elevatorStatus.ElevState = config.Idle
			} else {
				elevatorStatus.ElevState = config.Moving
			}
			setMotorDirection(elevatorStatus.Dir)
			return
		}
	}
}

func setMotorDirection(dir config.Directions) {
	if dir == config.MovingUp {
		elevio.SetMotorDirection(elevio.MD_Up)
	}
	if dir == config.MovingDown {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
}

<<<<<<< HEAD
func ElevStateMachine(ch config.FSMChannels, id int, sendOrder chan<- config.ElevatorOrder, newState chan<- map[string][config.NumElevators]config.ElevatorState,
	elevatorList *[config.NumElevators]config.ElevatorState, timerCh config.TimerChannels) {

	idAsString := strconv.Itoa(id)

	elevator := config.ElevatorState{
		Id:        id,
		Dir:       config.Stop,
		ElevState: config.Idle,
		Queue:     [config.NumFloors][config.NumBtns]bool{},
	}

	//go orderhandler.AddOrdersToQueue(newOrder, &elevator)
	initState(&elevator)

	//Stop elevator in the first floor that the elevators arrive in
	for {
		select {
		case floor := <-ch.Drv_floors:
			elevator.Floor = floor
			elevio.SetMotorDirection(elevio.MD_Stop)
			elevio.SetFloorIndicator(floor)
			break
		}
		break
	}

	elevatorList[id] = elevator
	newState <- map[string][config.NumElevators]config.ElevatorState{idAsString: *elevatorList}

	for {
		switch elevatorList[id].ElevState {
		case config.Idle:
			elevatorList[id].Dir = orderhandler.FindDirection(&elevatorList[id])
			if elevatorList[id].Dir != config.Stop {
				if elevatorList[id].Dir == config.MovingUp {
					elevio.SetMotorDirection(elevio.MD_Up)
				}
				if elevatorList[id].Dir == config.MovingDown {
					elevio.SetMotorDirection(elevio.MD_Down)
				}
				elevatorList[id].ElevState = config.Moving
			}
			if orderhandler.CheckOrderSameFLoor(&elevatorList[id], id) {
				elevatorList[id].ElevState = config.ArrivedAtFloor
			}
			if elevatorList[id].ElevState != config.Idle {
				go func() { newState <- map[string][config.NumElevators]config.ElevatorState{idAsString: *elevatorList} }()
			}

		case config.Moving:
			select {
			case floor := <-ch.Drv_floors:
				elevio.SetFloorIndicator(floor)
				elevatorList[id].Floor = floor
				if orderhandler.CheckIfArrived(floor, &elevatorList[id], id) {
					elevatorList[id].ElevState = config.ArrivedAtFloor
				}
				go func() { newState <- map[string][config.NumElevators]config.ElevatorState{idAsString: *elevatorList} }()
			}

		case config.ArrivedAtFloor:
			go func() { sendOrder <- config.ElevatorOrder{elevio.BT_HallUp, elevatorList[id].Floor, id, true} }()
			go func() { sendOrder <- config.ElevatorOrder{elevio.BT_HallDown, elevatorList[id].Floor, id, true} }()
			go timer.SetTimer(timerCh, config.Door)
			reachedFloor(timerCh.Open_door, &elevatorList[id])
			go func() { newState <- map[string][config.NumElevators]config.ElevatorState{idAsString: *elevatorList} }()
		}
	}
=======

func ElevStateMachine(ch config.FSMChannels, id int, sendOrder chan<- config.ElevatorOrder, sendState chan<- map[string][config.NumElevators]config.ElevatorState,
  elevatorList *[config.NumElevators]config.ElevatorState, timerCh config.TimerChannels) {

  idAsString := strconv.Itoa(id)
  idIndex := id - 1

  elevator := config.ElevatorState{
    Id:        id,
    Dir:       config.Stop,
    ElevState: config.Idle,
    Queue:     [config.NumFloors][config.NumBtns]bool{},
  }

  //go orderhandler.AddOrdersToQueue(newOrder, &elevator)
  initState(&elevator)

  //Stop elevator in the first floor that the elevators arrive in
  for {
    select{
    case floor := <- ch.Drv_floors:
      elevator.Floor = floor
      elevio.SetMotorDirection(elevio.MD_Stop)
      elevio.SetFloorIndicator(floor)
      break
    }
    break
  }

  elevatorList[idIndex] = elevator
  sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}

  for {
    switch elevatorList[idIndex].ElevState {
    case config.Idle:
        elevatorList[idIndex].Dir = orderhandler.FindDirection(&elevatorList[idIndex])
        if elevatorList[idIndex].Dir != config.Stop{
          if elevatorList[idIndex].Dir == config.MovingUp{
            elevio.SetMotorDirection(elevio.MD_Up)
          }
          if elevatorList[idIndex].Dir == config.MovingDown{
            elevio.SetMotorDirection(elevio.MD_Down)
          }
          elevatorList[idIndex].ElevState = config.Moving
        }
        if orderhandler.CheckOrderSameFLoor(&elevatorList[idIndex]){
          elevatorList[idIndex].ElevState = config.ArrivedAtFloor
        }
        if (elevatorList[idIndex].ElevState != config.Idle){
          sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
        }

    case config.Moving:
      select{
      case floor := <- ch.Drv_floors:
        elevio.SetFloorIndicator(floor)
        elevatorList[idIndex].Floor = floor
        if orderhandler.CheckIfArrived(floor, &elevatorList[idIndex]){
          elevatorList[idIndex].ElevState = config.ArrivedAtFloor
        }
        sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
      }

    case config.ArrivedAtFloor:
      sendOrder <- config.ElevatorOrder{elevio.BT_HallUp, elevatorList[idIndex].Floor, id, true}
      sendOrder <- config.ElevatorOrder{elevio.BT_HallDown, elevatorList[idIndex].Floor, id, true}
      sendOrder <- config.ElevatorOrder{elevio.BT_Cab, elevatorList[idIndex].Floor, id, true}
      go timer.SetTimer(timerCh, config.Door)
      reachedFloor(timerCh.Open_door, &elevatorList[idIndex])
      sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
    }
  }
>>>>>>> origin/development
}
