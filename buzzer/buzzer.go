package buzzer

import (
	"time"
)

// PwmChannel is an interface that abstracts PwmChannel functionality.
type PwmChannel interface {
	// Configure initializes the PWM channel.
	Configure() error
	// Top returns the maximum duty cycle value for the PWM channel.
	Top() uint32
	// SetPeriod sets the period for the PWM channel.
	SetPeriod(period uint64) error
	// SetDuty sets the duty cycle for the PWM channel.
	SetDuty(value uint32)
}

// NotePeriod represents the period of a note frequency in microseconds.
type NotePeriod int

const (
	C3  NotePeriod = 7634
	D3  NotePeriod = 6803
	E3  NotePeriod = 6061
	F3  NotePeriod = 5714
	G3  NotePeriod = 5102
	A3  NotePeriod = 4545
	B3  NotePeriod = 4049
	C4  NotePeriod = 3816 // 261 Hz
	D4  NotePeriod = 3401 // 294 Hz
	E4  NotePeriod = 3030 // 329 Hz
	F4  NotePeriod = 2865 // 349 Hz
	G4  NotePeriod = 2551 // 392 Hz
	A4  NotePeriod = 2272 // 440 Hz
	A4s NotePeriod = 2146
	B4  NotePeriod = 2028 // 493 Hz
	C5  NotePeriod = 1912 // 523 Hz
	D5  NotePeriod = 1706
	D5s NotePeriod = 1608
	E5  NotePeriod = 1517 // 659 Hz
	F5  NotePeriod = 1433 // 698 Hz
	G5s NotePeriod = 2408 // 784 Hz
	A5  NotePeriod = 1136
	A5s NotePeriod = 1073
	B5  NotePeriod = 1012
	C6  NotePeriod = 955

	Silence NotePeriod = 0
)

// Note represents a musical note with a specific frequency (period) and duration.
type Note struct {
	Period   NotePeriod // in microseconds
	Duration time.Duration
}

// Buzzer represents a buzzer that can play musical notes using a PWM channel.
type Buzzer struct {
	ch PwmChannel
}

// New creates a new Buzzer instance with the given PwmChannel.
func New(pwm PwmChannel) *Buzzer {
	return &Buzzer{
		ch: pwm,
	}
}

// Configure initializes the buzzer's PWM channel.
func (b *Buzzer) Configure() error {
	return b.ch.Configure()
}

// Tone plays a single note on the buzzer.
func (b *Buzzer) Tone(note Note) error {
	if note.Period == Silence {
		b.ch.SetDuty(0)
		time.Sleep(note.Duration)
		return nil
	}
	err := b.ch.SetPeriod(uint64(note.Period * 1000))
	if err != nil {
		return err
	}
	b.ch.SetDuty(b.ch.Top() / 2)
	time.Sleep(note.Duration)
	b.ch.SetDuty(0)
	return nil
}

// PlayMelody plays a sequence of notes (melody) on the buzzer.
func (b *Buzzer) PlayMelody(melody []Note) {
	for _, note := range melody {
		b.Tone(note)
	}
}
