package orderhandler

import "../elevio"
import "../config"

func CheckNewOrder(reciever chan<- config.ElevOrder, sender <-chan elevio.ButtonEvent, id int){
	for{
		select{
		case a := <- sender:
			order_floor := a.Floor
			button_type := a.Button
			executingElevator := id
			isDone := false;
			reciever <- config.ElevOrder{order_floor, button_type, executingElevator, isDone}
		}
	}
}
