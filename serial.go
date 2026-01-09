package main

import (
	"log"

	"go.bug.st/serial"
)

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
	defer func() {
		err := port.Close()
		if err != nil {
			log.Printf("erro ao fechar porta serial %s: %v\n", device, err)
			return
		}
	}()

	log.Printf("porta serial %s aberta\n", device)

	buf := make([]byte, 4096)
	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Printf("erro na leitura da serial: %v\n", err)
			return
		}
		if n > 5 {
			msg := string(buf[:n])
			err := MavControl.ParseRaw(msg)
			if err != nil {
				log.Println("Error reading sensor data: ", err.Error())
			}
		}
	}
}
