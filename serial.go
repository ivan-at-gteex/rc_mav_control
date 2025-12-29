package main

import (
	"log"

	"go.bug.st/serial"
)

type Joystick struct {
	X int
	Y int
}

// ReadSerial abre a porta serial indicada por "device" com baud padrão 115200
// e imprime continuamente todos os dados recebidos até o processo encerrar.
func ReadSerial(baud int, device string) {
	if device == "" {
		log.Println("porta serial não especificada")
		return
	}

	mode := &serial.Mode{BaudRate: baud}
	port, err := serial.Open(device, mode)
	if err != nil {
		log.Printf("erro ao abrir porta serial %s: %v\n", device, err)
		return
	}
	defer port.Close()

	log.Printf("porta serial %s aberta\n", device)

	buf := make([]byte, 4096)
	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Printf("erro na leitura da serial: %v\n", err)
			return
		}
		if n > 0 {
			// imprime exatamente os bytes recebidos (sem adicionar nova linha extra)
			log.Printf(string(buf[:n]))
		}
	}
}
