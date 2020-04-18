package orderhandler

import "../config"
import "../elevio"
import "time"
import "strconv"
//import "fmt"

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


//Making sure that the reconnecting elevator has the right state so that it can execute pre existing cab orders, as well as
//turning on lights
func syncElev(id int, tempElev config.ElevatorState, elevatorList *[config.NumElevators]config.ElevatorState){
	time.Sleep(3 * time.Second)
	for floor := 0; floor < config.NumFloors; floor++{
		if tempElev.Queue[floor][elevio.BT_Cab]{
			elevio.SetButtonLamp(elevio.BT_Cab, floor, tempElev.Queue[floor][elevio.BT_Cab])
			elevatorList[id].Queue[floor][elevio.BT_Cab] = true
		}
	}
	return
}

//Transfers the hall orders of the lost elevator to the best suited elevator on the network
func transferHallOrders(lostElevator config.ElevatorState, elevatorList *[config.NumElevators]config.ElevatorState, activeElevators *[config.NumElevators]bool,
	sendOrder chan<- config.ElevatorOrder, sendState chan<- map[string][config.NumElevators]config.ElevatorState, id int){
		idAsString := strconv.Itoa(id)
		lostElevatorIndex := lostElevator.Id-1
		for floor := 0; floor < config.NumFloors; floor++{
			for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{
				if lostElevator.Queue[floor][button]{
					elevatorList[lostElevatorIndex].Queue[floor][button] = false
					newExecutingElevator := costCalculator(floor, button, elevatorList, activeElevators, id)
					if newExecutingElevator == id{
						elevatorList[id-1].Queue[floor][button] = true
					}
					sendOrder <- config.ElevatorOrder{button, floor, newExecutingElevator, false}
					sendState <- map[string][config.NumElevators]config.ElevatorState{idAsString:*elevatorList}
				}
			}
		}
}

 func checkIfOthersAreActive( activeElevators *[config.NumElevators]bool , id int) bool{
	 for i :=0; i<len(activeElevators); i++{
		 if activeElevators[i] && i!=id{
			 return false
		 }
	 }
	 return true
 }

 
 func turnOnHallLightsWhenReconnectingToNetwork(sendersElevatorQueue [config.NumFloors][config.NumBtns]bool){
	 for floor := 0; floor < config.NumFloors; floor++{
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{
			if sendersElevatorQueue[floor][button]{
				elevio.SetButtonLamp(button, floor, true)
			}
		}
	 }
 }