package orderhandler

import "../elevio"
import "../config"
import "fmt"

func CheckNewOrder(reciever chan<- config.ElevatorOrder, sender <-chan elevio.ButtonEvent, id string){
	for{
		select{
		case a := <- sender:
			order_floor := a.Floor
			button_type := a.Button
			executingElevator := id
			isDone := false;
			fmt.Println("order CheckNewOrder")
			reciever <- config.ElevatorOrder{button_type, order_floor, executingElevator, isDone}
		}
	}
}


// func AddOrdersToQueue(sender <-chan config.ElevatorOrder, elevator *config.ElevatorState) {
// 	for{
// 		select{
// 		case newOrder := <- sender:
// 			order_type := newOrder.Button
// 			order_floor := newOrder.Floor
// 			elevio.SetButtonLamp(order_type, order_floor, true)
// 			elevator.Queue[order_floor][order_type] = true
// 		}
// 	}
// }

func OrderHandler(elevatorMap map[string]config.ElevatorState, activeElevators map[string]bool, buttonCh <-chan elevio.ButtonEvent){
		for{
			select{
			case pressedButton := <- buttonCh:
				button_type = pressedButton.Button
				order_floor = pressedButton.Floor
				best_elevator = CostCalculator(order_floor, button_type, elevatorMap, activeElevators)
				isDone := false
			}
		}

}
