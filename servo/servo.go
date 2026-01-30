package servo

type Servo180 interface {
	// SetAngle sets the angle of the servo in degrees (0-180)
	SetAngle(angle int) error
}

type Servo360 interface {
	// SetSpeed sets the speed of the servo in percentage (-100 to 100)
	SetSpeed(speed int) error
}
