package networkmod

import (
	"fmt"
	"../config"
	//"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.



func RecieveData(id string, ch config.NetworkChannels) {
	// Our id can be anything. Here we pass it on the command line, using
	//  `go run main.go -id=our_id`

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	// go peers.Transmitter(12346, id, peerTxEnable)
	// go peers.Receiver(12346, peerUpdateCh)

	// We make channels for sending and receiving our custom data types
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.




	// The example message. We just send one of these every second.

	/*go func() {
		helloMsg := HelloMsg{"Hello from " + id, 0}
		for {
			helloMsg.Iter++
			helloTx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()*/

	fmt.Println("Started")
	for {
		select {
		case p := <-ch.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)


		case a := <-ch.RecieveCh:
			fmt.Printf("Received: %#v\n", a)
		}
	}
}
