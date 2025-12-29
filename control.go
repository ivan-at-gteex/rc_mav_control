package main

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"sync"
)

const r = `\*(?P<input1>\d{1,4})\|(?P<input2>\d{1,4})\|(?P<input3>\d{1,4})\|(?P<input4>\d{1,4})\*`
const sensorIdx = 0.488400488

type Control struct {
	Joystick [2]Joystick
	Buttons  [12]bool
	mu       sync.Mutex
}

type Joystick struct {
	X int16
	Y int16
}

func (c *Control) Reset() {
	c.Joystick[0].X = 0
	c.Joystick[0].Y = 0
	c.Joystick[1].X = 0
	c.Joystick[1].Y = 0
}

func (c *Control) GetR() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return int16(math.Ceil(float64(c.Joystick[0].X)*sensorIdx) - 1000)
}

func (c *Control) GetZ() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return int16(math.Ceil(float64(c.Joystick[0].Y)*sensorIdx) - 1000)
}

func (c *Control) GetX() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return int16(math.Ceil(float64(c.Joystick[1].X)*sensorIdx) - 1000)
}

func (c *Control) GetY() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return int16(math.Ceil(float64(c.Joystick[1].Y)*sensorIdx) - 1000)
}

func (c *Control) ParseRaw(b []byte) error {

	InputRegex := regexp.MustCompile(r)
	subs := InputRegex.FindStringSubmatch(string(b))
	if len(subs) != 5 {
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

	c.mu.Lock()
	defer c.mu.Unlock()
	c.Joystick[0].X = int16(v1)
	c.Joystick[0].Y = int16(v2)
	c.Joystick[1].Y = int16(v3)
	c.Joystick[1].X = int16(v4)

	//log.Printf("Parsed frame: %d|%d|%d|%d", v1, v2, v3, v4)
	return nil
}
