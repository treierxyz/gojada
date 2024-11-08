package main

import (
	"github.com/treierxyz/gojada/devices"
	"github.com/treierxyz/gojada/devices/x55"
)

var Devices = map[string]devices.SerialDevice{
	"x55": x55.Device,
}
