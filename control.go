package main

type Control struct {
	Row      int
	Yaw      int
	Pitch    int
	Joystick [2]Joystick
	Buttons  [12]bool
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
