package networkmod

import (
	"./network/bcast"
	"./network/peers"
	"fmt"
	"time"
	"../config"
)


// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.
// type HelloMsg struct {
// 	Message string
// 	Iter    int
// 	Matix   [3][4]int
// }




// The example message. We just send one of these every second.


// sender <-chan config.ElevatorOrder

func SendData(id string) {

	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`


	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	peerUpdateCh := make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, id, peerTxEnable)
	go peers.Receiver(15647, peerUpdateCh)

	// We make channels for sending and receiving our custom data types

	transmitt:= make(chan config.ElevatorState)
	recieve := make(chan config.ElevatorState)
	
	// transmit_order := make(chan config.ElevatorOrder)
	// recieve_order := make(chan config.ElevatorOrder)
	
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.


	go bcast.Transmitter(16569, transmitt)
	go bcast.Receiver(16569, recieve)
	
	go func() {
		a := [4][3] bool{
			{false, false, false},
			{false, false, false},
			{false, false, false},
			{false, false, false}}
		helloMsg := config.ElevatorState{id, 1, 1, -1, a}
		//helloMsg := ElevatorState{"hello from "+ id, 0,  a}
		//helloMsg:={id, 0}

		for {
			//helloMsg.Iter++
			transmitt <- helloMsg
			time.Sleep(1 * time.Second)//må endre tid, sikkert sende meldinger mye oftere
		}
	}()

	fmt.Println("Started")
	for {
		select {
		case p := <-peerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
			if len(p.Lost) !=0 { //hvis du mister netverk med en heis
				fmt.Printf("oh no, you lost a elevator")
			}

		//case a := <-helloRx:
		//	fmt.Printf("Received: %#v\n", a)
		}
}
}
