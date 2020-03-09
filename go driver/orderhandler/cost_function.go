package orderhandler
//
// func costCalculator(order ElevOrder, elevs [NumElevators]ElevatorState, elevs_online [NumElevators]bool) int {
// 	if order.Button == ButtonUp {
// 		return elevs.ID
// 	}
// 	minCost := Inf(1)
// 	bestElevator := elevs.ID
// 	for i := 0; i < NumElevators; i++ { //iterating through all elevators
// 		if elevs_online[i] == false {
// 			continue //if elevator offline, skip to next iteration
// 		}
// 		cost := order.Floor - elevs[i].Floor
//
// 		if cost == 0 && elevs[i].ElevState == Idle {
// 			bestElevator = elevs[i].ID
// 			return bestElevator
// 		}
// 		if cost > 0 && elevs[i].Dir == MD_UP {
// 			cost = cost
// 		}
// 		if cost < 0 && elevs[i].Dir == MD_Down {
// 			//burde det ikke være cost--, for så å endre tegn helt til slutt?
// 			cost = cost
// 		}
// 		if cost > 0 && elevs[i].Dir == MD_Down {
// 			cost = cost + 2
// 		}
// 		if cost < 0 && elevs[i].Dir == MD_Up {
// 			//burde det ikke være cost = cost -3, -||-
// 			cost = cost + 2
// 		}
// 		if elevs[i].Dir == MD_Stop {
// 			cost = cost + 1
// 		}
// 		if cost == 0 && elevs[i].Dir != MD_Stop {
// 			cost = cost + 3
// 		}
// 		//burde det ikke være en if(cost<0){cost = -cost} her?
// 		if cost < minCost {
// 			minCost = cost
// 			bestElevator = elevs[i].ID
// 		}
// 		return bestElevator
// 	}
// }
