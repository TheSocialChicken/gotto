// Package servo defines interfaces for controlling different types of servos.
package servo

// Servo180 represents a servo with 180 degrees of rotation.
type Servo180 interface {
	// SetAngle sets the angle of the servo in degrees (0-180)
	SetAngle(angle int) error
}

// Servo360 defines the interface for a continuous rotation servo.
type Servo360 interface {
	// SetSpeed sets the speed of the servo in percentage (-100 to 100)
	SetSpeed(speed int) error
}
