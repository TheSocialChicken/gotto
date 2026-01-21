package ninja

type Servo180 interface {
	SetAngle(angle int) error
}

type Servo360 interface {
	SetSpeed(speed int) error
}
