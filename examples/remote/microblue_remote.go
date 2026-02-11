package main

import (
	"github.com/HattoriHanzo031/gotto/buzzer"
	"github.com/HattoriHanzo031/gotto/remote"
	"tinygo.org/x/bluetooth"
)

var (
	serviceUUID = bluetooth.ServiceUUIDNordicUART
	rxUUID      = bluetooth.CharacteristicUUIDUARTRX
)

type microBlue struct {
	ch chan byte
}

func NewMicroBlue() *microBlue {
	return &microBlue{
		ch: make(chan byte, 64),
	}
}

func (r *microBlue) ReadCommand() remote.Command {
	command := remote.Command{}
	state := byte(0)
	id := make([]byte, 0, 5)
	argIndex := 0

loop:
	for c := range r.ch {
		switch {
		case c == 1:
			state = 1
			id = id[:0]
		case c == 2:
			state = 2
			argIndex = 0
		case c == 3:
			break loop
		case state == 1:
			if len(id) >= 5 {
				continue
			}
			id = append(id, c)
		case state == 2:
			if c == ',' {
				if argIndex >= 2 {
					continue
				}
				argIndex++
				continue
			}
			command.Args[argIndex] = command.Args[argIndex]*10 + int(c-'0')
		}
	}

	switch string(id) {
	case "d1":
		command.Op = remote.OpRoll
		// map from 0..1023 to -100..100
		command.Args[0] = (command.Args[0]*200)/1023 - 100
		command.Args[1] = (command.Args[1]*200)/1023 - 100

	case "mo":
		command.Op = remote.OpSetMode
	case "tl":
		command.Op = remote.OpTiltLeft
	case "tr":
		command.Op = remote.OpTiltRight
	case "bz":
		command.Op = remote.OpBuzzerTone
		command.Args[1] = 500 * command.Args[0]
		command.Args[0] = int(buzzer.B3)
	case "rs":
		command.Op = remote.OpRightLegSpin
	case "ls":
		command.Op = remote.OpLeftLegSpin
	case "wa":
		if command.Args[0] == 0 {
			command.Op = remote.OpHome
		} else {
			command.Op = remote.OpWave
		}
	case "c0", "c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8", "c9":
		if command.Args[0] == 0 {
			command.Op = remote.OpHome
		} else {
			command.Op = remote.OpCustom
			command.Args[0] = int(id[1] - '0')
		}
	}
	return command
}

func (remote *microBlue) Start() error {
	adapter := bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		println("Error enabling BLE stack:", err.Error())
		return err
	}
	adv := adapter.DefaultAdvertisement()
	err := adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "GOtto",
		ServiceUUIDs: []bluetooth.UUID{serviceUUID},
	})
	if err != nil {
		println("Error configuring advertisement:", err.Error())
		return err
	}

	err = adv.Start()
	if err != nil {
		println("Error starting advertisement:", err.Error())
		return err
	}

	var rxChar bluetooth.Characteristic
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &rxChar,
				UUID:   rxUUID,
				Flags:  bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					for _, c := range value {
						select {
						case remote.ch <- c:
						default:
						}
					}
				},
			},
		},
	}))

	return nil
}
