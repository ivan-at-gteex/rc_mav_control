package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/ardupilotmega"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint.
// 2) print selected incoming messages.

var MavControl Control

func main() {

	c := context.Background()
	ctx, cancel := signal.NotifyContext(c, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

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

	MavControl.Init()

	go ReadEvents(node)
	go ReadSerial(460800, "/dev/ttyUSB0")

	log.Println("Program running")

	<-ctx.Done()
	cancel()
	log.Println("Received shutdown signal, shutting down.")
}

func ReadEvents(node *gomavlib.Node) {
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			fmt.Println("Received message, sending keyboard input")
			err := node.WriteMessageTo(frm.Channel, &ardupilotmega.MessageManualControl{
				Target:            frm.SystemID(),
				X:                 MavControl.GetX(),
				Y:                 MavControl.GetY(),
				Z:                 MavControl.GetZ(),
				R:                 MavControl.GetR(),
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
