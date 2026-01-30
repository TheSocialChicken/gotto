package servo

import "tinygo.org/x/drivers/servo"

type servo180 struct {
	servo  servo.Servo
	usLow  int
	usHigh int
}

func (s servo180) SetAngle(angle int) error {
	return s.servo.SetAngleWithMicroseconds(angle, s.usLow, s.usHigh)
}

func New180(s servo.Servo, usLow, usHigh int) servo180 {
	return servo180{
		servo:  s,
		usLow:  usLow,
		usHigh: usHigh,
	}
}

type servo360 struct {
	servo180
}

func (s servo360) SetSpeed(speed int) error {
	angle := speed + 100        // map -100..100 to 0..200
	angle = (angle * 180) / 200 // map 0..200 to 0..180
	return s.servo.SetAngleWithMicroseconds(angle, s.usLow, s.usHigh)
}

func New360(s servo.Servo, usLow, usHigh int) servo360 {
	return servo360{
		servo180: servo180{
			servo:  s,
			usLow:  usLow,
			usHigh: usHigh,
		},
	}
}
