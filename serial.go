package main

import (
	"time"

	"go.bug.st/serial"
)

var mode *serial.Mode = &serial.Mode{
	BaudRate: 19200,
}

func ConnectSerial(path string) serial.Port {
	port, err := serial.Open(path, mode)
	if err != nil {
		panic(err)
	}
	err = port.SetReadTimeout(50 * time.Millisecond)
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
