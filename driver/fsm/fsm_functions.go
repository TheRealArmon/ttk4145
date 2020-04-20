package fsm 

import (
	"../elevio"
	cf "../config"
	oh "../orderhandler"
  )

func initState(elevator *cf.ElevatorState, DrvFloors chan int, id int) {
	elevio.SetDoorOpenLamp(false)
	elevio.SetMotorDirection(elevio.MD_Down)
	  for {
	  select{
	  case floor := <- DrvFloors:
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetFloorIndicator(floor)
		elevator.Floor = floor
		elevator.Id = id
		for i := 0; i < cf.NumFloors; i++ {
		  for j := elevio.BT_HallUp; j < cf.NumBtns; j++ {
			elevio.SetButtonLamp(j, i, false)
		  }
		}
		return
	  }
	}
  }


func reachedFloor(door_timer <-chan bool, elevatorStatus *cf.ElevatorState) {
	oh.SwitchOffButtonLight(elevatorStatus.Floor)
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	elevatorStatus.Dir = cf.Stop
	for {
	  select{
	  case <- door_timer:
		elevio.SetDoorOpenLamp(false)
		if checkOrderSameFLoor(elevatorStatus){
		  elevatorStatus.State = cf.ArrivedAtFloor
		  return
		}
		elevatorStatus.Dir = findDirection(elevatorStatus)
		if (elevatorStatus.Dir == cf.Stop){
		  elevatorStatus.State = cf.Idle
		}else{elevatorStatus.State = cf.Moving}
		setMotorDirection(elevatorStatus.Dir)
		return
	  }
	}
  }

func findDirection(elevatorState *cf.ElevatorState) cf.Directions{
	switch elevatorState.Dir {
	case cf.Stop:
		if checkOrdersAbove(elevatorState){
			return cf.MovingUp
		}
		if checkOrdersBelow(elevatorState){
			return cf.MovingDown
		}else{return cf.Stop}
	
	case cf.MovingUp:
		if checkOrdersAbove(elevatorState){
			return cf.MovingUp
			}
		if checkOrdersBelow(elevatorState){
			return cf.MovingDown
		} else{return cf.Stop}
	
	case cf.MovingDown:
		if checkOrdersBelow(elevatorState){
			return cf.MovingDown
		}
		if checkOrdersAbove(elevatorState){
			return cf.MovingUp
		} else {return cf.Stop}
	}
	return cf.Stop
}

func checkOrdersAbove(elevatorState *cf.ElevatorState) bool{
	for floor := elevatorState.Floor + 1; floor < cf.NumFloors; floor++{
		for button := 0; button < cf.NumBtns; button++{
			if elevatorState.Queue[floor][button]{
				return true
			}
		}
	}
	return false
}

func checkOrdersBelow(elevatorState *cf.ElevatorState) bool{
	for floor := elevatorState.Floor - 1; floor > -1; floor--{
		for button := 0; button < cf.NumBtns; button++{
			if elevatorState.Queue[floor][button]{
				return true
			}
		}
	}
	return false
}


func checkIfArrived(floor int, elevatorState *cf.ElevatorState) bool{
	switch elevatorState.Dir {
	case cf.MovingUp:
		if elevatorState.Queue[floor][elevio.BT_Cab] || elevatorState.Queue[floor][elevio.BT_HallUp] || !checkOrdersAbove(elevatorState){
			return true
		}
	case cf.MovingDown:
		if elevatorState.Queue[floor][elevio.BT_Cab] || elevatorState.Queue[floor][elevio.BT_HallDown] || !checkOrdersBelow(elevatorState){
			return true
		}
	}
	return false
}
  
func checkOrderSameFLoor(elevatorState *cf.ElevatorState) bool{
	floor := elevatorState.Floor
	if elevatorState.Queue[floor][elevio.BT_Cab] ||
		elevatorState.Queue[floor][elevio.BT_HallUp] ||
			elevatorState.Queue[floor][elevio.BT_HallDown]{
			return true
	}
	return false
}

func clearOrderQueue(floor int, elevatorState *cf.ElevatorState){
	elevatorState.Queue[floor][elevio.BT_HallUp] = false
	elevatorState.Queue[floor][elevio.BT_Cab] = false
	elevatorState.Queue[floor][elevio.BT_HallDown] = false
}
  
func setMotorDirection(dir cf.Directions) {
	if dir == cf.MovingUp {
		elevio.SetMotorDirection(elevio.MD_Up)
	}
	if dir == cf.MovingDown {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
}

func findOrderButton(floor int, elevatorState *cf.ElevatorState) elevio.ButtonType{
	for button := elevio.BT_HallUp; button < cf.NumBtns; button++{
		if elevatorState.Queue[floor][button]{
			return button
		}
	}
	return elevio.BT_Cab
}



  