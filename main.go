package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint.
// 2) print selected incoming messages.

func main() {
	// create a node which communicates with a serial endpoint
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointUDPClient{Address: "192.168.2.15:8090"},
		},
		Dialect:     ardupilotmega.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 10,
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	go ReadEvents(node)

	fmt.Println("Aqui")

	for {
		fmt.Println("Aqui dentro (lá ele)")
		time.Sleep(10)
	}
}

func ReadEvents(node *gomavlib.Node) {
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			fmt.Println("Received message, sending keyboard input")
			err := node.WriteMessageTo(frm.Channel, &ardupilotmega.MessageManualControl{
				Target:            frm.SystemID(),
				X:                 50,
				Y:                 50,
				Z:                 50,
				R:                 0,
				Buttons:           0,
				Buttons2:          0,
				EnabledExtensions: 0,
				S:                 0,
				T:                 0,
				Aux1:              0,
				Aux2:              0,
				Aux3:              0,
				Aux4:              0,
				Aux5:              0,
				Aux6:              0,
			})
			if err != nil {
				log.Println("error writing frame:", err)
			}
		}
	}
	return
}
