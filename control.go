package main

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"sync"
)

const r = `\*(\d{1,4})\|(\d{1,4})\|(\d{1,4})\|(\d{1,4})\|(\d{10})\*`

//const r = `\*(?P<input1>\d{1,4})\|(?P<input2>\d{1,4})\|(?P<input3>\d{1,4})\|(?P<input4>\d{1,4})\|(?P<input5>\d{10})\*`

const centerRange = 10
const centerCounter = 1000

type Control struct {
	Joystick [2]Joystick
	Buttons  [10]bool
	mu       sync.Mutex
}

type Joystick struct {
	X Axis
	Y Axis
}

type Axis struct {
	current   int16
	max       int16
	min       int16
	center    int16
	centerSet bool
	history   History
	scaleMin  int16
	scaleMax  int16
	mu        sync.Mutex
}

type History struct {
	values         map[int16]int32
	mostOftenValue int16
}

func (a *Axis) Set(v int16) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.current = v

	if a.centerSet == false {
		a.history.values[a.current]++
		if a.history.values[a.current] >= a.history.values[a.history.mostOftenValue] {
			a.history.mostOftenValue = a.current
		}

		if a.history.values[a.current] >= centerCounter {
			a.centerSet = true
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

	scaled := int16(math.Ceil(float64(a.current-a.GetZero()) * a.GetScaleIndex()))

	if scaled > (centerRange*-1) && scaled < centerRange {
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
	return a.history.mostOftenValue
}

func (a *Axis) GetScaleIndex() float64 {
	valueRange := a.max - a.min
	scaleRange := a.scaleMax - a.scaleMin
	return float64(scaleRange) / float64(valueRange)
}

func (a *Axis) Init(scaleMin int16, scaleMax int16) {
	a.history.values = make(map[int16]int32)
	a.history.mostOftenValue = 0
	a.current = 0
	a.center = 0
	a.min = 100
	a.max = 3900
	a.scaleMin = scaleMin
	a.scaleMax = scaleMax
	a.centerSet = false
}

func (c *Control) Init() {
	c.Joystick[0].X.Init(-1000, 1000)
	c.Joystick[0].Y.Init(0, 1000)
	c.Joystick[1].X.Init(-1000, 1000)
	c.Joystick[1].Y.Init(-1000, 1000)
}

func (c *Control) GetR() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Joystick[0].X.GetScaled()
}

func (c *Control) GetZ() int16 {
	c.mu.Lock()
	defer c.mu.Unlock()

	scaleIndex := c.Joystick[0].Y.GetScaleIndex()
	return int16(math.Ceil(float64(c.Joystick[0].Y.Get())*scaleIndex)) - 100
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

func (c *Control) IsButtonPressed(button int) bool {
	return c.Buttons[button]
}

func (c *Control) ParseRaw(text string) error {

	InputRegex := regexp.MustCompile(r)
	subs := InputRegex.FindStringSubmatch(text)
	if len(subs) != 6 {
		return errors.New(fmt.Sprintf("invalid input, expected 5 values separated by |, got %d", len(subs)))
	}

	v1, err := strconv.Atoi(subs[1])
	if err != nil {
		return errors.Join(errors.New("invalid input for joystick 0, x axis"), errors.New(subs[1]), err)
	}
	v2, err := strconv.Atoi(subs[2])
	if err != nil {
		return errors.Join(errors.New("invalid input for joystick 0, Y axis"), errors.New(subs[2]), err)
	}
	v3, err := strconv.Atoi(subs[3])
	if err != nil {
		return errors.Join(errors.New("invalid input for joystick 1, y axis"), errors.New(subs[3]), err)
	}
	v4, err := strconv.Atoi(subs[4])
	if err != nil {
		return errors.Join(errors.New("invalid input for joystick 1, x axis"), errors.New(subs[4]), err)
	}

	for k, v := range subs[5] {
		if v == '1' {
			c.Buttons[k] = true
		} else {
			c.Buttons[k] = false
		}
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
