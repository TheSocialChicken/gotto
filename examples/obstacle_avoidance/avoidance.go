package main

import (
	"machine"
	"math/rand"
	"time"

	ninja "github.com/HattoriHanzo031/otto_ninja"
	"tinygo.org/x/drivers/hcsr04"
	"tinygo.org/x/drivers/servo"
)

var (
	pwmFoot = machine.PWM2
	pwmLeg  = machine.PWM1

	rFoot = machine.GP4
	lFoot = machine.GP5
	rLeg  = machine.GP2
	lLeg  = machine.GP3

	usTrig = machine.GP6
	usEcho = machine.GP7
)

type servo180 servo.Servo

func (s servo180) SetAngle(angle int) error {
	return servo.Servo(s).SetAngleWithMicroseconds(angle, 450, 2550)
}

type servo360 servo.Servo

func (s servo360) SetSpeed(speed int) error {
	angle := speed + 100        // map -100..100 to 0..200
	angle = (angle * 180) / 200 // map 0..200 to 0..180
	return servo.Servo(s).SetAngleWithMicroseconds(angle, 450, 2550)
}

func main() {
	time.Sleep(3 * time.Second)

	legArr := must(servo.NewArray(pwmLeg))
	footArr := must(servo.NewArray(pwmFoot))

	llServo := servo180(must(legArr.Add(lLeg)))
	rlServo := servo180(must(legArr.Add(rLeg)))
	lfServo := servo360(must(footArr.Add(lFoot)))
	rfServo := servo360(must(footArr.Add(rFoot)))

	n := ninja.New(rlServo, llServo, rfServo, lfServo)
	trim := ninja.Trim{
		TiltAngle:         0,
		LeftStepDuration:  0,
		RightStepDuration: 150 * time.Millisecond,
		LfSpeed:           0,
		RfSpeed:           0,
		LlAngle:           20,
		RlAngle:           12,
	}

	n.Trim(trim)

	us := hcsr04.New(usTrig, usEcho)
	us.Configure()

	for {
		// Walk 2 minutes with obstacle avoidance
		start := time.Now()
		n.Mode(ninja.ModeWalk)
		for time.Since(start) < 2*time.Minute {
			dist := us.ReadDistance()
			if dist != 0 && dist < 150 {
				n.Walk(-1)
				n.RightLegSpin(20, 100*time.Duration(rand.Intn(9)+6)*time.Millisecond)
				continue
			}

			n.Walk(1)
		}

		// Roll 2 minutes with obstacle avoidance
		start = time.Now()
		n.Mode(ninja.ModeRoll)
		n.Roll(50, 0)
		for time.Since(start) < 2*time.Minute {
			dist := us.ReadDistance()
			if dist != 0 && dist < 150 {
				n.RollStop()
				time.Sleep(500 * time.Millisecond)
				n.Roll(-50, 0)
				time.Sleep(500 * time.Millisecond)
				n.Roll(0, 50)
				rand.Intn(10)
				time.Sleep(100 * time.Duration(rand.Intn(9)+1) * time.Millisecond)
				n.Roll(50, 0)
			}
		}
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
