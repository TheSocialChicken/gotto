package main

import (
	"math/rand/v2"
	"time"

	"github.com/HattoriHanzo031/gotto/ninja"
	"tinygo.org/x/drivers/hcsr04"
)

func obstacleAvoidanceWalkFn(us hcsr04.Device) ninja.CustomCommand {
	return func(n *ninja.Ninja) error {
		start := time.Now()
		n.Mode(ninja.ModeWalk)
		for time.Since(start) < time.Minute {
			dist := us.ReadDistance()
			if dist != 0 && dist < 150 {
				n.Walk(-1)
				n.RightLegSpin(20, 100*time.Duration(rand.IntN(9)+6)*time.Millisecond)
				continue
			}

			n.Walk(1)
		}
		return nil
	}
}

func obstacleAvoidanceRollFn(us hcsr04.Device) ninja.CustomCommand {
	return func(n *ninja.Ninja) error {
		start := time.Now()
		n.Mode(ninja.ModeRoll)
		n.Roll(50, 0)
		for time.Since(start) < time.Minute {
			dist := us.ReadDistance()
			if dist != 0 && dist < 150 {
				n.RollStop()
				time.Sleep(500 * time.Millisecond)
				n.Roll(-50, 0)
				time.Sleep(500 * time.Millisecond)
				n.Roll(0, 50)
				time.Sleep(100 * time.Duration(rand.IntN(9)+1) * time.Millisecond)
				n.Roll(50, 0)
			}
		}
		return nil
	}
}
