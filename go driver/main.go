package main

import "./elevio"
import "./fsm"
import "./config"
//import "./orderhandler"




func main(){
    elevio.Init("localhost:15657", config.NumFloors)
    fsm.ElevStateMachine()

}
