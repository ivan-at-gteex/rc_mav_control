package main

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"sync"
)

const r = `\*(?P<input1>\d{1,4})\|(?P<input2>\d{1,4})\|(?P<input3>\d{1,4})\|(?P<input4>\d{1,4})\*`
const zeroRange = 10
const sensorIdx = 0.488400488

type Control struct {
	Joystick [2]Joystick
	Buttons  [12]bool
	mu       sync.Mutex
}

type Joystick struct {
	X Axis
	Y Axis
}

type Axis struct {
	current  int16
	max      int16
	min      int16
	zero     int16
	history  History
	ScaleMin int16
	ScaleMax int16
}

type History struct {
	values       map[int16]int16
	currentValue int16
}

func (a *Axis) Set(v int16) {
	a.current = v

	a.history.values[a.current]++
	if a.history.values[a.current] >= a.history.values[a.history.currentValue] {
		a.history.currentValue = a.current
	}

	if a.current <= a.min {
		a.min = a.current
		a.zero = (a.max - a.min) / 2
	}
	if a.current >= a.max {
		a.max = a.current
	}
}

func (a *Axis) Get() int16 {
	return a.current

}

func (a *Axis) GetScaled() int16 {
	return a.current

}

func (a *Axis) GetZero() int16 {
	return a.history.values[a.history.currentValue]
}

func (a *Axis) Init() {
	a.history.values = make(map[int16]int16)
	a.history.currentValue = 0
	a.current = 0
	a.zero = 0
	a.min = 1000
	a.max = -1000
	a.ScaleMin = -1000
	a.ScaleMax = 1000
}

func (c *Control) Init() {
	MavControl.Joystick[0].X.Init()
	MavControl.Joystick[0].Y.Init()
	MavControl.Joystick[1].X.Init()
	MavControl.Joystick[1].Y.Init()
}

func (c *Control) GetR() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return int16(math.Ceil(float64(c.Joystick[0].X.Get())*sensorIdx) - 1000)
}

func (c *Control) GetZ() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return int16(math.Ceil(float64(c.Joystick[0].Y.Get())*sensorIdx) - 1000)
}

func (c *Control) GetX() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return int16(math.Ceil(float64(c.Joystick[1].X.Get())*sensorIdx) - 1000)
}

func (c *Control) GetY() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()

	return int16(math.Ceil(float64(c.Joystick[1].Y.Get())*sensorIdx) - 1000)
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
		return errors.Join(errors.New("invalid input"), errors.New(subs[2]), err)
	}
	v3, err := strconv.Atoi(subs[3])
	if err != nil {
		return errors.Join(errors.New("invalid input"), errors.New(subs[3]), err)
	}
	v4, err := strconv.Atoi(subs[4])
	if err != nil {
		return errors.Join(errors.New("invalid input"), errors.New(subs[4]), err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.Joystick[0].X.Set(int16(v1))
	c.Joystick[0].Y.Set(int16(v2))
	c.Joystick[1].Y.Set(int16(v3))
	c.Joystick[1].X.Set(int16(v4))

	//log.Printf("Parsed frame: %d|%d|%d|%d", v1, v2, v3, v4)
	return nil
}
