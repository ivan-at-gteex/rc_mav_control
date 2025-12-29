package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sys/unix"
)

// KeyboardValue mantém o valor acumulado pelas teclas pressionadas.
// Seta para cima incrementa, seta para baixo decrementa.
var KeyboardValue int

func ReadKeyboard() {

	r := bufio.NewReader(os.Stdin)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return
		}

		fmt.Println("Keyboard input:", b)

		if b == 0x1b { // ESC
			// Possível sequência de escape: ESC [ A / ESC [ B
			b2, _ := r.ReadByte()
			if b2 != '[' {
				continue
			}
			b3, _ := r.ReadByte()
			switch b3 {
			case 'A': // seta para cima
				KeyboardValue++
			case 'B': // seta para baixo
				KeyboardValue--
			}
			continue
		}
		// Ignora outras teclas; somente setas afetam o valor
	}
}

func SetupKeyboard() {

	fd := int(os.Stdin.Fd())
	// Coloca o terminal em modo raw para ler teclas imediatamente (Unix)
	var oldState *unix.Termios
	if isTerminal(fd) {
		st, err := getTermios(fd)
		if err == nil {
			oldState = st
			raw := *st
			makeRaw(&raw)
			_ = setTermios(fd, &raw)
			defer func() {
				_ = setTermios(fd, oldState)
			}()
		}
	}

	// Garante restauração do terminal ao receber SIGINT/SIGTERM
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigC)
	go func() {
		<-sigC
		if oldState != nil {
			_ = setTermios(fd, oldState)
		}
		os.Exit(0)
	}()

	r := bufio.NewReader(os.Stdin)
	for {
		b, err := r.ReadByte()
		if err != nil {
			return
		}

		if b == 0x1b { // ESC
			// Possível sequência de escape: ESC [ A / ESC [ B
			b2, _ := r.ReadByte()
			if b2 != '[' {
				continue
			}
			b3, _ := r.ReadByte()
			switch b3 {
			case 'A': // seta para cima
				KeyboardValue++
			case 'B': // seta para baixo
				KeyboardValue--
			}
			continue
		}
		// Ignora outras teclas; somente setas afetam o valor
	}
}

// Helpers de terminal (Unix)
func isTerminal(fd int) bool {
	_, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	return err == nil
}

func getTermios(fd int) (*unix.Termios, error) {
	return unix.IoctlGetTermios(fd, unix.TCGETS)
}

func setTermios(fd int, t *unix.Termios) error {
	return unix.IoctlSetTermios(fd, unix.TCSETS, t)
}

func makeRaw(t *unix.Termios) {
	// baseado em cfmakeraw
	t.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	t.Oflag &^= unix.OPOST
	t.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	t.Cflag &^= unix.CSIZE | unix.PARENB
	t.Cflag |= unix.CS8
	t.Cc[unix.VMIN] = 1
	t.Cc[unix.VTIME] = 0
}
