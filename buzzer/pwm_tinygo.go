// TinyGo implementation of PWM channel for buzzer functionality.
package buzzer

import "machine"

// PWM is an interface that abstracts the functionality of tinygo PWM.
type PWM interface {
	Configure(config machine.PWMConfig) error
	Channel(pin machine.Pin) (channel uint8, err error)
	Top() uint32
	Set(channel uint8, value uint32)
}

type pwmChannel struct {
	pin     machine.Pin
	pwm     PWM
	channel uint8
}

// NewPwmChannel creates a new PwmChannel instance using the provided PWM and pin.
func NewPwmChannel(pwm PWM, pin machine.Pin) *pwmChannel {
	return &pwmChannel{
		pwm: pwm,
		pin: pin,
	}
}

// Configure initializes the PWM channel for the buzzer.
func (p *pwmChannel) Configure() error {
	channel, err := p.pwm.Channel(p.pin)
	if err != nil {
		return err
	}
	p.channel = channel

	return p.pwm.Configure(machine.PWMConfig{
		Period: 20000, // default period 20ms
	})
}

// Top returns the maximum duty cycle value for the PWM channel.
func (p *pwmChannel) Top() uint32 {
	return p.pwm.Top()
}

// SetDuty sets the duty cycle for the PWM channel.
func (p *pwmChannel) SetDuty(value uint32) {
	p.pwm.Set(p.channel, value)
}

// SetPeriod sets the period for the PWM channel.
func (p *pwmChannel) SetPeriod(period uint64) error {
	return p.pwm.Configure(machine.PWMConfig{
		Period: period,
	})
}
