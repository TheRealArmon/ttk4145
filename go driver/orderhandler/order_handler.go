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

//TODO lag en funksjon som