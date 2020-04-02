package orderhandler

import "../config"
import "../elevio"
import "fmt"
//import "sync"

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
				return id
			}
			if cost > 0 && elevator.Dir == config.MovingDown {
				cost += 3
			}
			if cost < 0{
				//fmt.Println("cost er negativ for heis nr: ", id)
				cost = -cost
				if elevator.Dir == config.MovingUp{
					//fmt.Println("Plusser på 3 for heis nr:", id)
					cost += 3
				}
			}
			if elevator.ElevState == config.ArrivedAtFloor {
				cost++
			}
			if cost == 0 && elevator.Dir != config.Stop {
				cost += 4
			}
			/*fmt.Printf("%v ", id)
			fmt.Printf("has cost %v ", cost)
			fmt.Printf("with dircetion, %v", elevator.Dir)
			fmt.Println(" at floor", floor)*/
			if cost < minCost {
				minCost = cost
				bestElevator = id
				
			}
		}
	fmt.Println("Dette er den beste heisen: ", bestElevator)
	fmt.Println("")
	return bestElevator
}
