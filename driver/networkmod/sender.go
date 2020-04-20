package networkmod


import (
	"time"
	"../config"
)

//Sends data 10 times with a frequency of 20 per second
func SendData(networkCh config.NetworkChannels, orderCh config.OrderChannels){
	interval := 15 * time.Millisecond
	for {
		select{
		case orderMsg := <- orderCh.SendOrder:
			for i := 0; i < 10; i++{
				networkCh.TransmittOrderCh <- orderMsg
				time.Sleep(interval)
			}

		case stateMsg := <- orderCh.SendState:
			for i := 0; i < 10; i++{
				networkCh.TransmittStateCh <- stateMsg
				time.Sleep(interval)
			}
		}	
	}	
}