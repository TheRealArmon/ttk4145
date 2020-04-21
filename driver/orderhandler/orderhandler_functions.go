package orderhandler


import (
	cf "../config"
 	"../elevio"
	"time"
	"strconv"
	"fmt"
)

func costCalculator(floor int, button_type elevio.ButtonType, elevatorList *[cf.NumElevators]cf.ElevatorState,
	 activeElevators *[cf.NumElevators]bool, id int) int {
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
		if cost == 0 && elevator.State == cf.Idle {
			return elevator.Id
		}
		if cost > 0 && elevator.Dir == cf.MovingDown {
			cost += 3
		}
		if cost < 0{
			cost = -cost
			if elevator.Dir == cf.MovingUp{
				cost += 3
			}
		}
		if elevator.State == cf.ArrivedAtFloor {
			cost++
		}
		if cost == 0 && elevator.Dir != cf.Stop {
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
func syncReconnectedElev(id int, tempElev cf.ElevatorState, elevatorList *[cf.NumElevators]cf.ElevatorState){
	time.Sleep(3 * time.Second)
	for floor := 0; floor < cf.NumFloors; floor++{
		if tempElev.Queue[floor][elevio.BT_Cab]{
			elevio.SetButtonLamp(elevio.BT_Cab, floor, tempElev.Queue[floor][elevio.BT_Cab])
			elevatorList[id].Queue[floor][elevio.BT_Cab] = true
		}
	}
	return
}

//Transfers the hall orders of the lost elevator to the best suited elevator on the network
func transferHallOrders(lostElevator cf.ElevatorState, elevatorList *[cf.NumElevators]cf.ElevatorState, activeElevators *[cf.NumElevators]bool,
	orderCh cf.OrderChannels, id int){
		idAsString := strconv.Itoa(id)
		lostElevatorIndex := lostElevator.Id-1
		for floor := 0; floor < cf.NumFloors; floor++{
			for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{
				if lostElevator.Queue[floor][button]{
					elevatorList[lostElevatorIndex].Queue[floor][button] = false
					newExecutingElevator := costCalculator(floor, button, elevatorList, activeElevators, id)
					if newExecutingElevator == id{
						elevatorList[id-1].Queue[floor][button] = true
					}
					orderCh.SendOrder <- cf.ElevatorOrder{button, floor, newExecutingElevator, false}
					orderCh.SendState <- map[string][cf.NumElevators]cf.ElevatorState{idAsString:*elevatorList}
				}
			}
		}
}

 func checkIfOthersAreActive( activeElevators *[cf.NumElevators]bool , id int) bool{
	 for i :=0; i<len(activeElevators); i++{
		 if activeElevators[i] && i!=id{
			 return true
		 }
	 }
	 return false
 }

 //If an elevator is executing hall orders when an another elevator connects, the connected elevator turn on lights
 func turnOnHallLightsWhenReconnectingToNetwork(sendersElevatorQueue [cf.NumFloors][cf.NumBtns]bool){
	 time.Sleep(1 * time.Second)
	 for floor := 0; floor < cf.NumFloors; floor++{
		for button := elevio.BT_HallUp; button < elevio.BT_Cab; button++{
			if sendersElevatorQueue[floor][button]{
				elevio.SetButtonLamp(button, floor, true)
			}
		}
	 }
 }

 func SwitchOffButtonLight(floor int){
	for button := elevio.BT_HallUp; button < cf.NumBtns; button++{
	  elevio.SetButtonLamp(button, floor, false)
	}
  }