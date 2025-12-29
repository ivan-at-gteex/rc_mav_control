package main

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"sync"
)

const r = `\*(?P<input1>\d{1,4})\|(?P<input2>\d{1,4})\|(?P<input3>\d{1,4})\|(?P<input4>\d{1,4})\*`

type Control struct {
	Row      int
	Yaw      int
	Pitch    int
	Joystick [2]Joystick
	Buttons  [12]bool
	mu       sync.Mutex
}

type Joystick struct {
	X int
	Y int
}

func (c *Control) Reset() {
	c.Row = 0
	c.Yaw = 0
	c.Pitch = 0
}

func (c *Control) GetRow() int {
	return 500
}

func (c *Control) ParseRaw(b []byte) error {

	InputRegex := regexp.MustCompile(r)
	subs := InputRegex.FindStringSubmatch(string(b))
	if len(subs) != 5 { // whole match + 4 groups
		// Should not happen, but break to avoid infinite loop
		return errors.New("invalid input")
	}

	v1, err := strconv.Atoi(subs[1])
	if err != nil {
		return errors.Join(errors.New("invalid input"), errors.New(subs[1]), err)
	}
	v2, err := strconv.Atoi(subs[2])
	if err != nil {
		return errors.Join(errors.New("invalid input"), errors.New(subs[1]), err)
	}
	v3, err := strconv.Atoi(subs[3])
	if err != nil {
		return errors.Join(errors.New("invalid input"), errors.New(subs[1]), err)
	}
	v4, err := strconv.Atoi(subs[4])
	if err != nil {
		return errors.Join(errors.New("invalid input"), errors.New(subs[1]), err)
	}

	c.Joystick[0].X = v1
	c.Joystick[0].Y = v2
	c.Joystick[1].X = v3
	c.Joystick[1].Y = v4

	log.Printf("Parsed frame: %d|%d|%d|%d", v1, v2, v3, v4)
	return nil
}
