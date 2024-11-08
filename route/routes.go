package route

import (
	"context"
	"encoding/hex"

	"github.com/danielgtaylor/huma/v2"
	"github.com/treierxyz/hk-projector-api/devices"
	"github.com/treierxyz/hk-projector-api/serialhelper"
	"go.bug.st/serial"
)

type RawSendOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello world" doc:"Device reply"`
	}
}

func CreateRoutes(api huma.API, port serial.Port, device devices.SerialDevice) {
	for _, cmd := range device.Commands() { // iterate over commands
		huma.Register(api, cmd.Operation, func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
			serialhelper.WriteSerial(port, cmd.Data)
			resp := &RawSendOutput{}
			result := serialhelper.ReadSerial(port)
			resp.Body.Message = hex.EncodeToString(result)
			return resp, nil
		}) // register each command
	}
}
