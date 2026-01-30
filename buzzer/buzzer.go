package buzzer

import (
	"machine"
	"time"
)

type PWM interface {
	Configure(config machine.PWMConfig) error
	Channel(pin machine.Pin) (channel uint8, err error)
	Top() uint32
	Set(channel uint8, value uint32)
}

type NotePeriod int

const (
	C3  = 7634
	D3  = 6803
	E3  = 6061
	F3  = 5714
	G3  = 5102
	A3  = 4545
	B3  = 4049
	C4  = 3816 // 261 Hz
	D4  = 3401 // 294 Hz
	E4  = 3030 // 329 Hz
	F4  = 2865 // 349 Hz
	G4  = 2551 // 392 Hz
	A4  = 2272 // 440 Hz
	A4s = 2146
	B4  = 2028 // 493 Hz
	C5  = 1912 // 523 Hz
	D5  = 1706
	D5s = 1608
	E5  = 1517 // 659 Hz
	F5  = 1433 // 698 Hz
	G5s = 2408 // 784 Hz
	A5  = 1136
	A5s = 1073
	B5  = 1012
	C6  = 955

	Silence = 0
)

type Note struct {
	Period   NotePeriod // in microseconds
	Duration time.Duration
}

type Buzzer struct {
	tempo   NotePeriod
	pwm     PWM
	pin     machine.Pin
	channel uint8
}

func New(pwm PWM, pin machine.Pin) *Buzzer {
	return &Buzzer{
		pwm:   pwm,
		pin:   pin,
		tempo: 1000,
	}
}

func (b *Buzzer) Configure() error {
	err := b.pwm.Configure(machine.PWMConfig{
		Period: 1e9 / 2000, // 2kHz
	})
	if err != nil {
		return err
	}
	ch, err := b.pwm.Channel(b.pin)
	if err != nil {
		return err
	}
	b.channel = ch
	return nil
}

func (b *Buzzer) SetTempo(tempo NotePeriod) {
	b.tempo = tempo
}

func (b *Buzzer) Tone(note Note) error {
	if note.Period == Silence {
		b.pwm.Set(b.channel, 0)
		time.Sleep(note.Duration)
		return nil
	}
	err := b.pwm.Configure(machine.PWMConfig{Period: uint64(note.Period * b.tempo)})
	if err != nil {
		return err
	}
	b.pwm.Set(b.channel, b.pwm.Top()/4)
	time.Sleep(note.Duration)
	b.pwm.Set(b.channel, 0)
	return nil
}

func (b *Buzzer) PlayMelody(melody []Note) {
	for _, note := range melody {
		b.Tone(note)
	}
}
