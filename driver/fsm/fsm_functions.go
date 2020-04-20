package fsm 

import "../elevio"
import "../config"
import "../orderhandler"

func initState(elevator *config.ElevatorState, DrvFloors chan int, id int) {
	elevio.SetDoorOpenLamp(false)
	elevio.SetMotorDirection(elevio.MD_Down)
	  for {
	  select{
	  case floor := <- DrvFloors:
		elevio.SetMotorDirection(elevio.MD_Stop)
		elevio.SetFloorIndicator(floor)
		elevator.Floor = floor
		elevator.Id = id
		for i := 0; i < config.NumFloors; i++ {
		  for j := elevio.BT_HallUp; j < config.NumBtns; j++ {
			elevio.SetButtonLamp(j, i, false)
		  }
		}
		return
	  }
	}
  }


func reachedFloor(door_timer <-chan bool, elevatorStatus *config.ElevatorState) {
	elevatorStatus.Dir = config.Stop
	orderhandler.SwitchOffButtonLight(elevatorStatus.Floor)
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	for {
	  select{
	  case <- door_timer:
		elevio.SetDoorOpenLamp(false)
		if orderhandler.CheckOrderSameFLoor(elevatorStatus){
		  elevatorStatus.State = config.ArrivedAtFloor
		  return
		}
		elevatorStatus.Dir = orderhandler.FindDirection(elevatorStatus)
		if (elevatorStatus.Dir == config.Stop){
		  elevatorStatus.State = config.Idle
		}else{elevatorStatus.State = config.Moving}
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
