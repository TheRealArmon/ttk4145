package orderhandler

import "../config"
import "../elevio"

func costCalculator(floor int, button_type elevio.ButtonType, elevatorList *[config.NumElevators]config.ElevatorState,
	 activeElevators *[config.NumElevators]bool, id int) int {
	if button_type == elevio.BT_Cab {
		return id
	}
	minCost := 1000
	bestElevator := id
	cost := 0
	idIndex := id - 1
	for id, elevator := range elevatorList{ //iterating through all elevators
		if !activeElevators[id] {
			continue //if elevator offline, skip to next iteration
		}
		cost = floor - elevator.Floor
		if (activeElevators[idIndex]){
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
	}
	return bestElevator
}
