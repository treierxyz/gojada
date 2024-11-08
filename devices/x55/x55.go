package x55

import (
	"bytes"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/treierxyz/hk-projector-api/devices"
	"go.bug.st/serial"
)

type X55Device struct {
	commands []devices.SerialDeviceCommand
}

type X55CmdData struct {
	crc   [2]byte
	ctype [2]byte
	code  [2]byte
}

var Header = [5]byte{0xbe, 0xef, 0x03, 0x06, 0x00}

var Actions = map[string][]byte{
	"set":       {0x01, 0x00},
	"get":       {0x02, 0x00},
	"increment": {0x04, 0x00},
	"decrement": {0x05, 0x00},
	"execute":   {0x06, 0x00},
}

func (X55Device) Settings() devices.SerialDeviceSettings {
	return devices.SerialDeviceSettings{
		Mode: serial.Mode{
			BaudRate: 19200,
		},
		ReadTimeout: 100 * time.Millisecond,
	}
}

func (device X55Device) Commands() []devices.SerialDeviceCommand {
	for i, cmd := range device.commands {
		if device.commands[i].Operation.Path == "" {
			device.commands[i].Operation.Path = "/" + strings.Replace(cmd.Name, "-", "/", -1)
		}
		device.commands[i].Operation.OperationID = cmd.Name
		device.commands[i].Operation.Method = methodFromAction(cmd.Data)
		device.commands[i].Operation.Description += "\n\n" + "Raw command: 0x" + hex.EncodeToString(cmd.Data)
	}
	return device.commands
}

func createBytes(action []byte, data X55CmdData) []byte {
	return append(append(append(append(Header[:], data.crc[:]...), action[:]...), data.ctype[:]...), data.code[:]...)
}

func methodFromAction(data []byte) string {
	for actionName, action := range Actions {
		if bytes.Equal(data[7:9], action) {
			switch actionName {
			case "set", "increment", "decrement", "execute":
				return http.MethodPost
			case "get":
				return http.MethodGet
			default:
				return ""
			}
		}
	}
	return ""
}

var Device = X55Device{
	[]devices.SerialDeviceCommand{
		{
			Name: "power-on",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0xba, 0xd2},
				ctype: [2]byte{0x00, 0x60},
				code:  [2]byte{0x01, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Turn power on",
				Description: "Turns the projector on",
				Tags:        []string{"Power", "Set"},
			},
		},
		{
			Name: "power-off",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x2a, 0xd3},
				ctype: [2]byte{0x00, 0x60},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Turn power off",
				Description: "Turns the projector off",
				Tags:        []string{"Power", "Set"},
			},
		},
		{
			Name: "power-get",
			Data: createBytes(Actions["get"], X55CmdData{
				crc:   [2]byte{0x19, 0xd3},
				ctype: [2]byte{0x00, 0x60},
				code:  [2]byte{0x00, 0x00},
			}),
			Responses: map[string][]byte{
				"off":      {0x00, 0x00},
				"on":       {0x01, 0x00},
				"cooldown": {0x02, 0x00},
			},
			Operation: huma.Operation{
				Summary:     "Get power state",
				Description: "Retrieves the projectors power state",
				Tags:        []string{"Power", "Get"},
			},
		},
		{
			Name: "input-rgb1",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0xfe, 0xd2},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Set input to RGB1",
				Description: "Sets the input to RGB1",
				Tags:        []string{"Input", "Set"},
			},
		},
		{
			Name: "input-rgb2",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x3e, 0xd0},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x04, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Set input to RGB2",
				Description: "Sets the input to RGB2",
				Tags:        []string{"Input", "Set"},
			},
		},
		{
			Name: "input-video",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x6e, 0xd3},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x01, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Set input to Video",
				Description: "Sets the input to Video",
				Tags:        []string{"Input", "Set"},
			},
		},
		{
			Name: "input-s-video",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x9e, 0xd3},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x02, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/input/s-video",
				Summary:     "Set input to S-Video",
				Description: "Sets the input to S-Video",
				Tags:        []string{"Input", "Set"},
			},
		},
		{
			Name: "input-component",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0xae, 0xd1},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x05, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Set input to Component",
				Description: "Sets the input to Component",
				Tags:        []string{"Input", "Set"},
			},
		},
		{
			Name: "input-get",
			Data: createBytes(Actions["get"], X55CmdData{
				crc:   [2]byte{0xcd, 0xd2},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x00, 0x00},
			}),
			Responses: map[string][]byte{
				"rgb1":      {0x00, 0x00},
				"rgb2":      {0x04, 0x00},
				"video":     {0x01, 0x00},
				"s-video":   {0x02, 0x00},
				"component": {0x05, 0x00},
			},
			Operation: huma.Operation{
				Summary:     "Get input",
				Description: "Retrieves the projectors current input",
				Tags:        []string{"Input", "Get"},
			},
		},
		{
			Name: "error-get",
			Data: createBytes(Actions["get"], X55CmdData{
				crc:   [2]byte{0xd9, 0xd8},
				ctype: [2]byte{0x20, 0x60},
				code:  [2]byte{0x00, 0x00},
			}),
			Responses: map[string][]byte{
				"normal":          {0x00, 0x00},
				"cover_error":     {0x01, 0x00},
				"fan_error":       {0x02, 0x00},
				"lamp_error":      {0x03, 0x00},
				"temp_error":      {0x04, 0x00},
				"air_flow_error":  {0x05, 0x00},
				"lamp_time_error": {0x06, 0x00},
				"cool_error":      {0x07, 0x00},
				"filter_error":    {0x08, 0x00},
			},
			Operation: huma.Operation{
				Summary:     "Get error status",
				Description: "Retrieves the projectors error status",
				Tags:        []string{"Error", "Get"},
			},
		},
		{
			Name: "brightness-get",
			Data: createBytes(Actions["get"], X55CmdData{
				crc:   [2]byte{0x89, 0xd2},
				ctype: [2]byte{0x03, 0x20},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Get brightness",
				Description: "Retrieves the projectors brightness",
				Tags:        []string{"Brightness", "Get"},
			},
		},
		{
			Name: "brightness-increment",
			Data: createBytes(Actions["increment"], X55CmdData{
				crc:   [2]byte{0xef, 0xd2},
				ctype: [2]byte{0x03, 0x20},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Increment brightness",
				Description: "Increments the projectors brightness",
				Tags:        []string{"Brightness", "Increment"},
			},
		},
		{
			Name: "brightness-decrement",
			Data: createBytes(Actions["decrement"], X55CmdData{
				crc:   [2]byte{0x3e, 0xd3},
				ctype: [2]byte{0x03, 0x20},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Decrement brightness",
				Description: "Decrements the projectors brightness",
				Tags:        []string{"Brightness", "Decrement"},
			},
		},
		{
			Name: "brightness-reset",
			Data: createBytes(Actions["execute"], X55CmdData{
				crc:   [2]byte{0x58, 0xd3},
				ctype: [2]byte{0x00, 0x70},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Reset brightness",
				Description: "Resets the projectors brightness",
				Tags:        []string{"Brightness", "Reset"},
			},
		},
		{
			Name: "gamma-default-1",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x07, 0xe9},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x20, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/gamma/default-1",
				Summary:     "Set gamma to \"Default 1\"",
				Description: "Sets the gamma to \"Default 1\"",
				Tags:        []string{"Gamma", "Set"},
			},
		},
		{
			Name: "gamma-default-2",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x97, 0xe8},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x21, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/gamma/default-2",
				Summary:     "Set gamma to \"Default 2\"",
				Description: "Sets the gamma to \"Default 2\"",
				Tags:        []string{"Gamma", "Set"},
			},
		},
		{
			Name: "gamma-default-3",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x67, 0xe8},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x22, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/gamma/default-3",
				Summary:     "Set gamma to \"Default 3\"",
				Description: "Sets the gamma to \"Default 3\"",
				Tags:        []string{"Gamma", "Set"},
			},
		},
		{
			Name: "gamma-custom-1",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x07, 0xfd},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x10, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/gamma/custom-1",
				Summary:     "Set gamma to \"Custom 1\"",
				Description: "Sets the gamma to \"Custom 1\"",
				Tags:        []string{"Gamma", "Set"},
			},
		},
		{
			Name: "gamma-custom-2",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x97, 0xfc},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x11, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/gamma/custom-2",
				Summary:     "Set gamma to \"Custom 2\"",
				Description: "Sets the gamma to \"Custom 2\"",
				Tags:        []string{"Gamma", "Set"},
			},
		},
		{
			Name: "gamma-custom-3",
			Data: createBytes(Actions["set"], X55CmdData{
				crc:   [2]byte{0x67, 0xfc},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x12, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/gamma/custom-3",
				Summary:     "Set gamma to \"Custom 3\"",
				Description: "Sets the gamma to \"Custom 3\"",
				Tags:        []string{"Gamma", "Set"},
			},
		},
		{
			Name: "gamma-get",
			Data: createBytes(Actions["get"], X55CmdData{
				crc:   [2]byte{0xf4, 0xf0},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Summary:     "Get gamma",
				Description: "Retrieves the projectors gamma",
				Tags:        []string{"Gamma", "Get"},
			},
		},
		{
			Name: "lamp-time-get",
			Data: createBytes(Actions["get"], X55CmdData{
				crc:   [2]byte{0xc2, 0xff},
				ctype: [2]byte{0x90, 0x10},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/lamp-time/get",
				Summary:     "Get lamp time",
				Description: "Retrieves the projectors lamp time",
				Tags:        []string{"Lamp time", "Get"},
			},
		},
		{
			Name: "lamp-time-reset",
			Data: createBytes(Actions["execute"], X55CmdData{
				crc:   [2]byte{0x58, 0xdc},
				ctype: [2]byte{0x30, 0x70},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/lamp-time/reset",
				Summary:     "Reset lamp time",
				Description: "Resets the projectors lamp time",
				Tags:        []string{"Lamp time", "Reset"},
			},
		},
		{
			Name: "filter-time-get",
			Data: createBytes(Actions["get"], X55CmdData{
				crc:   [2]byte{0xc2, 0xf0},
				ctype: [2]byte{0xa0, 0x10},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/filter-time/get",
				Summary:     "Get filter time",
				Description: "Retrieves the projectors filter time",
				Tags:        []string{"Filter time", "Get"},
			},
		},
		{
			Name: "filter-time-reset",
			Data: createBytes(Actions["execute"], X55CmdData{
				crc:   [2]byte{0x98, 0xc6},
				ctype: [2]byte{0x40, 0x70},
				code:  [2]byte{0x00, 0x00},
			}),
			Operation: huma.Operation{
				Path:        "/filter-time/reset",
				Summary:     "Reset filter time",
				Description: "Resets the projectors filter time",
				Tags:        []string{"Filter time", "Reset"},
			},
		},
	},
}
