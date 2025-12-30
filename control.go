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
const zeroCounter = 1000

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
	zeroSet  bool
	history  History
	scaleMin int16
	scaleMax int16
	mu       sync.Mutex
}

type History struct {
	values       map[int16]int32
	currentValue int16
}

func (a *Axis) Set(v int16) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.current = v

	if a.zeroSet == false {
		a.history.values[a.current]++
		if a.history.values[a.current] >= a.history.values[a.history.currentValue] {
			a.history.currentValue = a.current
		}

		if a.history.values[a.current] >= zeroCounter {
			a.zeroSet = true
		}
	}

	if a.current <= a.min {
		a.min = a.current
	}
	if a.current >= a.max {
		a.max = a.current
	}
}

func (a *Axis) Get() int16 {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.current

}

func (a *Axis) GetScaled() int16 {
	a.mu.Lock()
	defer a.mu.Unlock()

	valueRange := a.max - a.min
	scaleRange := a.scaleMax - a.scaleMin
	scaleIndex := float64(scaleRange) / float64(valueRange)
	scaled := int16(math.Ceil(float64(a.current-a.GetZero()) * scaleIndex))

	if scaled > (zeroRange*-1) && scaled < zeroRange {
		return 0
	}

	if scaled > a.scaleMax {
		return a.scaleMax
	}

	if scaled < a.scaleMin {
		return a.scaleMin
	}

	return scaled
}

func (a *Axis) GetZero() int16 {
	return a.history.currentValue
}

func (a *Axis) Init(scaleMin int16, scaleMax int16) {
	a.history.values = make(map[int16]int32)
	a.history.currentValue = 0
	a.current = 0
	a.zero = 0
	a.min = 500
	a.max = 3500
	a.scaleMin = scaleMin
	a.scaleMax = scaleMax
	a.zeroSet = false
}

func (c *Control) Init() {
	MavControl.Joystick[0].X.Init(-1000, 1000)
	MavControl.Joystick[0].Y.Init(-1000, 1000)
	MavControl.Joystick[1].X.Init(-1000, 1000)
	MavControl.Joystick[1].Y.Init(-1000, 1000)
}

func (c *Control) GetR() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Joystick[0].X.GetScaled()
}

func (c *Control) GetZ() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Joystick[0].Y.GetScaled()
}

func (c *Control) GetX() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Joystick[1].X.GetScaled()
}

func (c *Control) GetY() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Joystick[1].Y.GetScaled()
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
