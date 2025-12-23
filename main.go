// Package main contains an example.
package main

import (
	"log"
	"net"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/pion/transport/v2/udp"
)

// this example shows how to:
// 1) create a node which communicates with a serial endpoint.
// 2) print selected incoming messages.

func main() {
	// create a node which communicates with a serial endpoint
	node := &gomavlib.Node{
		Endpoints: []gomavlib.EndpointConf{
			gomavlib.EndpointSerial{
				Device: "/dev/ttyACM0",
				Baud:   57600,
			},
			gomavlib.EndpointCustomServer{
				Listen: func() (net.Listener, error) {
					addr, err := net.ResolveUDPAddr("udp4", ":8090")
					if err != nil {
						return nil, err
					}

					return udp.Listen("udp4", addr)
				},
				Label: "udp",
			},
		},
		Dialect:     nil,         // do not use a dialect and do not attempt to decode messages (in a router it is preferable)
		OutVersion:  gomavlib.V2, // change to V1 if you're unable to communicate with the target
		OutSystemID: 10,
	}
	err := node.Initialize()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	// print incoming frames
	for evt := range node.Events() {
		if frm, ok := evt.(*gomavlib.EventFrame); ok {
			log.Printf("received: id=%d, %+v\n", frm.Message().GetID(), frm.Message())
			// route frame to every other channel
			err = node.WriteFrameExcept(frm.Channel, frm.Frame)
			if err != nil {
				log.Println("error writing frame:", err)
			}
		}
	}
}
