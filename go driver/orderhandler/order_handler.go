package orderhandler

import "../elevio"
import "../config"

func CheckNewOrder(reciever chan<- config.ElevatorOrder, sender <-chan elevio.ButtonEvent, id int){
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
}
