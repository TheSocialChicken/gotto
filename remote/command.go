// Package remote provides functionality to create and execute remote control commands for the Ninja robot.
package remote

import (
	"errors"
	"time"

	"github.com/HattoriHanzo031/gotto/buzzer"
	"github.com/HattoriHanzo031/gotto/ninja"
)

// Op represents the operation code for a remote control command.
type Op int

const (
	OpSetMode Op = iota
	OpWalk
	OpRoll
	OpTiltRight
	OpTiltLeft
	OpHome
	OpRightLegSpin
	OpLeftLegSpin
	OpBuzzerTone
	OpWave
	OpCustom
)

var (
	ErrUnknownCommand = errors.New("ninja: invalid command")
)

// Command represents a remote control command for the Ninja robot.
type Command struct {
	// Op specifies the operation to be performed (e.g., set mode, walk, roll, etc.).
	Op Op
	// Args contains the arguments for the command, such as speed, direction, etc.
	Args [2]int
}

// Execute performs the command on the given Ninja instance.
func (c *Command) Execute(n *ninja.Ninja) error {
	switch c.Op {
	case OpSetMode:
		switch c.Args[0] {
		case 0:
			return n.Mode(ninja.ModeWalk)
		case 1:
			return n.Mode(ninja.ModeRoll)
		}
	case OpHome:
		return n.Home()
	case OpTiltLeft:
		switch c.Args[0] {
		case 0:
			return n.Tilt(ninja.TiltReturnFromLeft)
		case 1:
			return n.Tilt(ninja.TiltLeft)
		}
	case OpTiltRight:
		switch c.Args[0] {
		case 0:
			return n.Tilt(ninja.TiltReturnFromRight)
		case 1:
			return n.Tilt(ninja.TiltRight)
		}
	case OpLeftLegSpin:
		if c.Args[0] == 1 {
			return n.StartLeftSpin(30)
		} else {
			return n.StopLeftSpin()
		}
	case OpRightLegSpin:
		if c.Args[0] == 1 {
			return n.StartRightSpin(30)
		} else {
			return n.StopRightSpin()
		}
	case OpWalk:
		return n.Walk(c.Args[0])
	case OpRoll:
		return n.Roll(c.Args[1], c.Args[0])
	case OpBuzzerTone:
		return n.BuzzerTone(buzzer.Note{
			Period:   buzzer.NotePeriod(c.Args[0]),
			Duration: time.Duration(c.Args[1]) * time.Millisecond,
		})
	case OpWave:
		return n.Wave()
	case OpCustom:
		return n.ExecuteCustomCommand(c.Args[0])
	default:
		return ErrUnknownCommand
	}
	return nil
}
