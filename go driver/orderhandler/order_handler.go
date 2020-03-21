package orderhandler

import "../elevio"
import "../config"
import "fmt"
 /*
 func CheckNewOrder(reciever chan<- config.ElevatorOrder, sender <-chan elevio.ButtonEvent, id string){
	for{
		select{
		case a := <- sender:
			order_floor := a.Floor
			button_type := a.Button
			executingElevator := id
			isDone := false;
			reciever <- config.ElevatorOrder{button_type, order_floor, executingElevator, isDone}
		}
	}
}

 func AddOrdersToQueue(sender <-chan config.ElevatorOrder, elevator *config.ElevatorState) {
	for{
		select{
 		case newOrder := <- sender:
 			order_type := newOrder.Button
 			order_floor := newOrder.Floor
 			elevio.SetButtonLamp(order_type, order_floor, true)
 			elevator.Queue[order_floor][order_type] = true
 		}
 	}
 } */

func OrderHandler(elevatorMap map[string]config.ElevatorState, activeElevators map[string]bool, buttonCh <-chan elevio.ButtonEvent, id string){
		for{
			select{
			case pressedButton := <- buttonCh:
				button_type := pressedButton.Button
				order_floor := pressedButton.Floor
				best_elevator := costCalculator(order_floor, button_type, elevatorMap, activeElevators, id)
				isDone := false
				fmt.Println(best_elevator)
				fmt.Println(isDone)
			}
		}

}

func RecievePeer(sender <-chan []string){
	for{
		select{
		case  <- sender:

		}
	}
}

func RecievedStateUpdateFromNetwork(state <-chan map[string]config.ElevatorState, elevatorMap map[string]config.ElevatorState){
	for{
		select{
		case newState := <- state:
			for id, elevatorState := range newState{
				if val, ok := elevatorMap[id]; ok {
					fmt.Println(val)
   					elevatorMap[id] = elevatorState	
				}
			}
			fmt.Println(len(elevatorMap))
		}
	}
}
