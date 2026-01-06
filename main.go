package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"rc_mavlink/config"
	"syscall"
	"time"

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

	cfg, errLoad := config.Load()
	if errLoad != nil {
		panic(errLoad)
	}

	serverAddr := cfg.MavLinkAddr + ":" + cfg.MavlinkPort

	// create a node which communicates with a serial endpoint
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointUDPClient{Address: serverAddr},
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
	go ReadSerial(cfg.SerialBaud, cfg.SerialPort)

	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			log.Println("Control Values: ",
				MavControl.Joystick[0].X.Get(),
				MavControl.Joystick[0].Y.Get(),
				MavControl.Joystick[1].X.Get(),
				MavControl.Joystick[1].Y.Get())

			log.Println("Control Scaled Values: ",
				MavControl.Joystick[0].X.GetScaled(),
				MavControl.Joystick[0].Y.GetScaled(),
				MavControl.Joystick[1].X.GetScaled(),
				MavControl.Joystick[1].Y.GetScaled())
		}
	}()

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
