package serialhelper

import (
	"github.com/treierxyz/hk-projector-api/devices"
	"go.bug.st/serial"
)

func ConnectSerial(path string, settings devices.SerialDeviceSettings) serial.Port {
	port, err := serial.Open(path, &settings.Mode)
	if err != nil {
		panic(err)
	}
	err = port.SetReadTimeout(settings.ReadTimeout)
	if err != nil {
		panic(err)
	}
	return port
}

func DisconnectSerial(port serial.Port) {
	err := port.Close()
	if err != nil {
		panic(err)
	}
}

func WriteSerial(port serial.Port, data []byte) int {
	n, err := port.Write(data)
	if err != nil {
		panic(err)
	}
	return n
}

// read serial until no more data is available
func ReadSerial(port serial.Port) []byte {
	var final []byte
	for {
		buf := make([]byte, 8)
		n, err := port.Read(buf)

		if err != nil {
			panic(err)
		}
		if n == 0 {
			// no more data available
			return final
		}

		final = append(final, buf[:n]...)
	}
}
