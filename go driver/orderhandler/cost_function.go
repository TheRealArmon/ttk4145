package orderhandler

import "../config"

func CostCalculator(floor int, button_type elevio.ButtonType, elevatorMap map[string]config.ElevatorState, activeElevators map[string]bool, id string) string {
	if button_type == elevio.BT_Cab {
		return id
	}
	minCost := Inf(1)
	bestElevator := id
  cost := 0
	for elevator := range elevatorMap{ //iterating through all elevators
		if !activeElevators[elevator] {
			continue //if elevator offline, skip to next iteration
		}
		cost = floor - elevatorMap[elevator].Floor

		if cost == 0 && elevatorMap[elevator].ElevState == config.Idle {
			return elevator
		}
		if cost > 0 && elevatorMap[elevator].Dir == config.MovingUp {
			cost = cost
		}
		if cost < 0 && elevatorMap[elevator].Dir == config.MovingDown {
			//burde det ikke være cost--, for så å endre tegn helt til slutt?
			cost = cost
		}
		if cost > 0 && elevatorMap[elevator].Dir == config.MovingDown {
			cost = cost + 2
		}
		if cost < 0 && elevatorMap[elevator].Dir == config.MovingUp {
			//burde det ikke være cost = cost -3, -||-
			cost = cost + 2
		}
		if elevatorMap[elevator].Dir == config.Stop {
			cost = cost + 1
		}
		if cost == 0 && elevatorMap[elevator].Dir != config.Stop {
			cost = cost + 3
		}
		//burde det ikke være en if(cost<0){cost = -cost} her?
		if cost < minCost {
			minCost = cost
			bestElevator = elevator
		}
		return bestElevator
	}
}
