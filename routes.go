package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"reflect"

	"github.com/danielgtaylor/huma/v2"
	"go.bug.st/serial"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var caserLower = cases.Lower(language.English)
var caserTitle = cases.Title(language.English)

func fieldsIntoSlice(val reflect.Value, action [2]byte) []byte {
	// create byte slice out of SerialTXCmdGeneric fields
	// because reflect.Value is unaddressable, we have to iterate over every byte and make a new slice

	if val.Type() == reflect.TypeOf(SerialTXCmdGeneric{}) { // check if value is type SerialTXCmdGeneric
		crc := make([]byte, 0)
		ctype := make([]byte, 0)
		code := make([]byte, 0)

		for i := 0; i < val.FieldByName("crc").Len(); i++ {
			crc = append(crc, byte(val.FieldByName("crc").Index(i).Uint()))
		}

		for i := 0; i < val.FieldByName("ctype").Len(); i++ {
			ctype = append(ctype, byte(val.FieldByName("ctype").Index(i).Uint()))
		}

		for i := 0; i < val.FieldByName("code").Len(); i++ {
			code = append(code, byte(val.FieldByName("code").Index(i).Uint()))
		}

		// create byte slice from byte slices
		return append(append(append(append(Header[:], crc...), action[:]...), ctype...), code...)
	} else {
		return nil
	}
}

func CreateRoutes(api huma.API, port serial.Port, cmds SerialTXCommands) {
	for cmdName, cmd := range cmds { // iterace over commands
		cmdVal := reflect.ValueOf(cmd)
		for i := 0; i < cmdVal.NumField(); i++ { // iterate over actions
			actionName := cmdVal.Type().Field(i).Name
			newpath := "/" + cmdName + "/" + actionName

			if cmdVal.Field(i).IsZero() { // skip empty actions
				continue
			}

			switch cmdVal.Type().Field(i).Type { // switch over action types
			case reflect.TypeOf(SerialTXCmdSet{}): // Set action
				for iter := cmdVal.Field(i).MapRange(); iter.Next(); { // iterate over set subcommands
					key := iter.Key()
					val := iter.Value()
					// fmt.Println("SerialTXCmdSet ## POST " + newpath + "/" + key.String())
					// fmt.Println("--- Value: ", val)

					action := Actions[actionName]

					// create byte slice from byte slices
					bytes := fieldsIntoSlice(val, action)

					huma.Register(api, huma.Operation{
						OperationID: caserLower.String(actionName + "-" + caserLower.String(cmdName) + "-" + key.String()),
						Method:      http.MethodPost,
						Path:        caserLower.String(newpath + "/" + key.String()),
						Summary:     "Set " + cmdName + " to " + key.String(),
						Description: "Sets the " + cmdName + " to " + key.String() + "\n\nCommand sent: 0x" + hex.EncodeToString(bytes),
						Tags:        []string{cmdName, actionName},
					}, func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
						WriteSerial(port, bytes)
						resp := &RawSendOutput{}
						result := ReadSerial(port)
						resp.Body.Message = hex.EncodeToString(result)
						return resp, nil
					})
				}
			case reflect.TypeOf(SerialTXCmdGet{}): // Get action
				// fmt.Println("SerialTXCmdGet ## GET " + newpath)
				// fmt.Println("--- Value: ", cmdVal.Field(i).FieldByName("SerialTXCmdGeneric"))

				action := Actions[cmdVal.Type().Field(i).Name]
				bytes := fieldsIntoSlice(cmdVal.Field(i).FieldByName("SerialTXCmdGeneric"), action)

				var routefunc func(ctx context.Context, input *struct{}) (*RawSendOutput, error)
				if ress := cmdVal.Field(i).FieldByName("responses"); ress.Len() > 0 {
					routefunc = func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
						WriteSerial(port, bytes)
						resp := &RawSendOutput{}
						result := ReadSerial(port)
						resp.Body.Message = hex.EncodeToString(result)
						return resp, nil
					}
				} else {
					routefunc = func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
						WriteSerial(port, bytes)
						resp := &RawSendOutput{}
						result := ReadSerial(port)
						resp.Body.Message = hex.EncodeToString(result)
						return resp, nil
					}
				}
				huma.Register(api, huma.Operation{
					OperationID: caserLower.String(actionName + "-" + cmdName),
					Method:      http.MethodGet,
					Path:        caserLower.String(newpath),
					Summary:     "Get " + cmdName + " value",
					Description: "Gets the value of " + cmdName + "\n\nCommand sent: 0x" + hex.EncodeToString(bytes),
					Tags:        []string{cmdName, actionName},
				}, routefunc)
			case reflect.TypeOf(SerialTXCmdGeneric{}): // Generic action (increment, decrement, execute)
				// fmt.Println("SerialTXCmdGeneric ## POST " + newpath)
				// fmt.Println("--- Value: ", cmdVal.Field(i))
				action := Actions[cmdVal.Type().Field(i).Name]
				bytes := fieldsIntoSlice(cmdVal.Field(i), action)
				var actualActionName string

				if actionName == "execute" && cmdName != "Auto-Adjust" {
					actualActionName = "Reset"
				} else {
					actualActionName = actionName
				}

				huma.Register(api, huma.Operation{
					OperationID: caserLower.String(actualActionName + "-" + caserLower.String(cmdName)),
					Method:      http.MethodPost,
					Path:        caserLower.String(newpath),
					Summary:     caserTitle.String(actualActionName) + "s " + cmdName,
					Description: caserTitle.String(actualActionName) + "s " + cmdName + "\n\nCommand sent: 0x" + hex.EncodeToString(bytes),
					Tags:        []string{cmdName, actualActionName},
				}, func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
					WriteSerial(port, bytes)
					resp := &RawSendOutput{}
					result := ReadSerial(port)
					resp.Body.Message = hex.EncodeToString(result)
					return resp, nil
				})
			default:
				fmt.Println("Unknown")
			}
		}
	}
}
