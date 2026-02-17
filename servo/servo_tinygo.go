// TinyGo servo wrapper for ninja servos
package servo

import "tinygo.org/x/drivers/servo"

// servo180 represents a servo with 180 degrees of rotation.
type servo180 struct {
	servo  servo.Servo
	usLow  int
	usHigh int
}

// New180 creates a new 180-degree servo with specified microsecond range.
func New180(s servo.Servo, usLow, usHigh int) servo180 {
	return servo180{
		servo:  s,
		usLow:  usLow,
		usHigh: usHigh,
	}
}

// SetAngle sets the angle of the servo in degrees (0-180)
func (s servo180) SetAngle(angle int) error {
	return s.servo.SetAngleWithMicroseconds(angle, s.usLow, s.usHigh)
}

// servo360 represents a continuous rotation servo.
type servo360 struct {
	servo180
}

// New360 creates a new 360-degree servo with specified microsecond range.
func New360(s servo.Servo, usLow, usHigh int) servo360 {
	return servo360{
		servo180: servo180{
			servo:  s,
			usLow:  usLow,
			usHigh: usHigh,
		},
	}
}

// SetSpeed sets the speed of the servo in percentage (-100 to 100)
func (s servo360) SetSpeed(speed int) error {
	angle := speed + 100        // map -100..100 to 0..200
	angle = (angle * 180) / 200 // map 0..200 to 0..180
	return s.SetAngle(angle)
}
