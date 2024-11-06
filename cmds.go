package main

type (
	SerialTXCommands = map[SerialTXCmdName]SerialTXActions
	SerialTXCmdName  = string
	SerialTXCmdSet   = map[string]SerialTXCmdGeneric
)

type SerialTXActions struct {
	set       SerialTXCmdSet
	get       SerialTXCmdGet
	increment SerialTXCmdGeneric
	decrement SerialTXCmdGeneric
	execute   SerialTXCmdGeneric
}

type SerialTXCmdGeneric struct {
	crc   [2]byte
	ctype [2]byte
	code  [2]byte
}

type SerialTXCmdGet struct {
	responses map[SerialTXCmdName][2]byte
	SerialTXCmdGeneric
}

var Header = [5]byte{0xbe, 0xef, 0x03, 0x06, 0x00}

var Actions = map[string][2]byte{
	"set":       {0x01, 0x00},
	"get":       {0x02, 0x00},
	"increment": {0x04, 0x00},
	"decrement": {0x05, 0x00},
	"execute":   {0x06, 0x00},
}

var Commands = SerialTXCommands{
	"Power": SerialTXActions{
		set: SerialTXCmdSet{
			"off": SerialTXCmdGeneric{
				crc:   [2]byte{0x2a, 0xd3},
				ctype: [2]byte{0x00, 0x60},
				code:  [2]byte{0x00, 0x00},
			},
			"on": SerialTXCmdGeneric{
				crc:   [2]byte{0xba, 0xd2},
				ctype: [2]byte{0x00, 0x60},
				code:  [2]byte{0x01, 0x00},
			},
		},
		get: SerialTXCmdGet{
			SerialTXCmdGeneric: SerialTXCmdGeneric{
				crc:   [2]byte{0x19, 0xd3},
				ctype: [2]byte{0x00, 0x60},
				code:  [2]byte{0x00, 0x00},
			},
			responses: map[SerialTXCmdName][2]byte{
				"off":      {0x00, 0x00},
				"on":       {0x01, 0x00},
				"cooldown": {0x02, 0x00},
			},
		},
	},
	"Input": SerialTXActions{
		set: SerialTXCmdSet{
			"RGB1": SerialTXCmdGeneric{
				crc:   [2]byte{0xfe, 0xd2},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x00, 0x00},
			},
			"RGB2": SerialTXCmdGeneric{
				crc:   [2]byte{0x3e, 0xd0},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x04, 0x00},
			},
			"Video": SerialTXCmdGeneric{
				crc:   [2]byte{0x6e, 0xd3},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x01, 0x00},
			},
			"S-Video": SerialTXCmdGeneric{
				crc:   [2]byte{0x9e, 0xd3},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x02, 0x00},
			},
			"Component": SerialTXCmdGeneric{
				crc:   [2]byte{0xae, 0xd1},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x05, 0x00},
			},
		},
		get: SerialTXCmdGet{
			SerialTXCmdGeneric: SerialTXCmdGeneric{
				crc:   [2]byte{0xcd, 0xd2},
				ctype: [2]byte{0x00, 0x20},
				code:  [2]byte{0x00, 0x00},
			},
			responses: map[SerialTXCmdName][2]byte{
				"rgb1":      {0x00, 0x00},
				"rgb2":      {0x04, 0x00},
				"video":     {0x01, 0x00},
				"s-video":   {0x02, 0x00},
				"component": {0x05, 0x00},
			},
		},
	},
	"Error-Status": SerialTXActions{
		get: SerialTXCmdGet{
			SerialTXCmdGeneric: SerialTXCmdGeneric{
				crc:   [2]byte{0xd9, 0xd8},
				ctype: [2]byte{0x20, 0x60},
				code:  [2]byte{0x00, 0x00},
			},
			responses: map[SerialTXCmdName][2]byte{
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
		},
	},
	"Brightness": SerialTXActions{
		get: SerialTXCmdGet{
			SerialTXCmdGeneric: SerialTXCmdGeneric{
				crc:   [2]byte{0x89, 0xd2},
				ctype: [2]byte{0x03, 0x20},
				code:  [2]byte{0x00, 0x00},
			},
			responses: make(map[SerialTXCmdName][2]byte),
		},
		increment: SerialTXCmdGeneric{
			crc:   [2]byte{0xef, 0xd2},
			ctype: [2]byte{0x03, 0x20},
			code:  [2]byte{0x00, 0x00},
		},
		decrement: SerialTXCmdGeneric{
			crc:   [2]byte{0x3e, 0xd3},
			ctype: [2]byte{0x03, 0x20},
			code:  [2]byte{0x00, 0x00},
		},
		execute: SerialTXCmdGeneric{
			crc:   [2]byte{0x58, 0xd3},
			ctype: [2]byte{0x00, 0x70},
			code:  [2]byte{0x00, 0x00},
		},
	},
	"Gamma": SerialTXActions{
		set: SerialTXCmdSet{
			"Default-1": SerialTXCmdGeneric{
				crc:   [2]byte{0x07, 0xe9},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x20, 0x00},
			},
			"Default-2": SerialTXCmdGeneric{
				crc:   [2]byte{0x97, 0xe8},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x21, 0x00},
			},
			"Default-3": SerialTXCmdGeneric{
				crc:   [2]byte{0x67, 0xe8},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x22, 0x00},
			},
			"Custom-1": SerialTXCmdGeneric{
				crc:   [2]byte{0x07, 0xfd},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x10, 0x00},
			},
			"Custom-2": SerialTXCmdGeneric{
				crc:   [2]byte{0x97, 0xfc},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x11, 0x00},
			},
			"Custom-3": SerialTXCmdGeneric{
				crc:   [2]byte{0x67, 0xfc},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x12, 0x00},
			},
		},
		get: SerialTXCmdGet{
			SerialTXCmdGeneric: SerialTXCmdGeneric{
				crc:   [2]byte{0xf4, 0xf0},
				ctype: [2]byte{0xa1, 0x30},
				code:  [2]byte{0x00, 0x00},
			},
			responses: make(map[SerialTXCmdName][2]byte),
		},
	},
	"Lamp-Time": SerialTXActions{
		get: SerialTXCmdGet{
			SerialTXCmdGeneric: SerialTXCmdGeneric{
				crc:   [2]byte{0xc2, 0xff},
				ctype: [2]byte{0x90, 0x10},
				code:  [2]byte{0x00, 0x00},
			},
			responses: make(map[SerialTXCmdName][2]byte),
		},
		execute: SerialTXCmdGeneric{
			crc:   [2]byte{0x58, 0xdc},
			ctype: [2]byte{0x30, 0x70},
			code:  [2]byte{0x00, 0x00},
		},
	},
	"Filter-Time": SerialTXActions{
		get: SerialTXCmdGet{
			SerialTXCmdGeneric: SerialTXCmdGeneric{
				crc:   [2]byte{0xc2, 0xf0},
				ctype: [2]byte{0xa0, 0x10},
				code:  [2]byte{0x00, 0x00},
			},
			responses: make(map[SerialTXCmdName][2]byte),
		},
		execute: SerialTXCmdGeneric{
			crc:   [2]byte{0x98, 0xc6},
			ctype: [2]byte{0x40, 0x70},
			code:  [2]byte{0x00, 0x00},
		},
	},
}
