package devices

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"go.bug.st/serial"
)

type SerialDevice interface {
	Settings() SerialDeviceSettings
	Commands() []SerialDeviceCommand
}

type SerialDeviceSettings struct {
	serial.Mode
	ReadTimeout time.Duration
}

type SerialDeviceCommand struct {
	Name      string
	Data      []byte
	Responses map[string][]byte
	huma.Operation
}

// TODO: Add handler function to struct, I can't figure that out yet
