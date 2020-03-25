package orderhandler

import "../elevio"
import "../config"
//import "sync"
//import "fmt"



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

func OrderHandler(buttonCh <-chan elevio.ButtonEvent, sendOrder chan<- config.ElevatorOrder, id int,
	elevatorList *[config.NumElevators]config.ElevatorState, activeElevators *[config.NumElevators]bool){
		for{
			select{
			case pressedButton := <- buttonCh:
				button_type := pressedButton.Button
				order_floor := pressedButton.Floor
				elevio.SetButtonLamp(button_type, order_floor, true)
				best_elevator := costCalculator(order_floor, button_type, elevatorList, activeElevators, id)
				isDone := false
				sendOrder <- config.ElevatorOrder{button_type, order_floor, best_elevator, isDone}
			}
		}
}


