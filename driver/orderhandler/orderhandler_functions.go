package orderhandler

import "../config"
import "../elevio"
import "time"

func costCalculator(floor int, button_type elevio.ButtonType, elevatorList *[config.NumElevators]config.ElevatorState,
	 activeElevators *[config.NumElevators]bool, id int) int {
	if button_type == elevio.BT_Cab {
		return id
	}
	minCost := 1000
	bestElevator := id
	cost := 0
	for id, elevator := range elevatorList{ //iterating through all elevators
		if !activeElevators[id] {
			continue //if elevator offline, skip to next iteration
		}
		cost = floor - elevator.Floor
		if cost == 0 && elevator.ElevState == config.Idle {
			return elevator.Id
		}
		if cost > 0 && elevator.Dir == config.MovingDown {
			cost += 3
		}
		if cost < 0{
			cost = -cost
			if elevator.Dir == config.MovingUp{
				cost += 3
			}
		}
		if elevator.ElevState == config.ArrivedAtFloor {
			cost++
		}
		if cost == 0 && elevator.Dir != config.Stop {
			cost += 4
		}
		if cost < minCost {
			minCost = cost
			bestElevator = elevator.Id

		}
	}
	return bestElevator
}

func checkCabQueue(elevatorState config.ElevatorState) bool {
	for floor := 0; floor < config.NumFloors; floor++ {
		if elevatorState.Queue[floor][elevio.BT_Cab] {
			return true
		}
	}
	return false
}


//Making sure that the reconnecting elevator has the right state so that it can execute pre existing cab orders, as well as
//turning on lights
func syncElev(id int, tempElev config.ElevatorState, elevatorList *[config.NumElevators]config.ElevatorState){
	if tempElev.Dir == -1 && tempElev.ElevState != config.ArrivedAtFloor{
		tempElev.Floor -= 1
	}
	for floor := 0; floor < config.NumFloors; floor++{
		elevio.SetButtonLamp(elevio.BT_Cab, floor, tempElev.Queue[floor][elevio.BT_Cab])
	}
	tempElev.ElevState = config.Idle
	tempElev.Dir = config.Stop
	time.Sleep(3 * time.Second)
	elevatorList[id] = tempElev
	return
}