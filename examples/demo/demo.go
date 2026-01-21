package main

import (
	"machine"

	ninja "github.com/HattoriHanzo031/otto_ninja"
	"tinygo.org/x/drivers/servo"
)

var (
	pwmLeg  = machine.PWM0
	pwmFoot = machine.PWM1

	rLeg  = machine.GP0
	lLeg  = machine.GP1
	rFoot = machine.GP2
	lFoot = machine.GP3
)

type servo180 servo.Servo

func (s servo180) SetAngle(angle int) error {
	return servo.Servo(s).SetAngleWithMicroseconds(angle, 500, 2500)
}

type servo360 servo.Servo

func (s servo360) SetSpeed(speed int) error {
	angle := speed + 100        // map -100..100 to 0..200
	angle = (angle * 180) / 200 // map 0..200 to 0..180
	return servo.Servo(s).SetAngleWithMicroseconds(angle, 500, 2500)
}

func main() {
	legArr := must(servo.NewArray(pwmLeg))
	footArr := must(servo.NewArray(pwmFoot))

	n := ninja.New(
		servo180(must(legArr.Add(rLeg))),
		servo180(must(legArr.Add(lLeg))),
		servo360(must(footArr.Add(rFoot))),
		servo360(must(footArr.Add(lFoot))))

	n.Mode(ninja.ModeRoll)
	n.RollStop()
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
