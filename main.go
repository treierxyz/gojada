package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Options for the CLI.
type Options struct {
	Port   int    `help:"Port to listen on" short:"p" default:"8888"`
	Serial string `help:"Path of emulated device serial" short:"E" default:"/dev/ttyUSB0"`
}

type RawSendInput struct {
	Message string `path:"message" maxLength:"300" example:"world" doc:"Message to send"`
}

type RawSendOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello world" doc:"Device reply"`
	}
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		fmt.Println("Connecting to serial...")
		port := ConnectSerial(options.Serial)
		fmt.Println("Connected successfully to serial ports")

		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("Häkkerikoda Projector API", "0.1.0"))

		huma.Post(api, "/power/set/on", func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
			WriteSerial(port, []byte{0xbe, 0xef, 0x03, 0x06, 0x00, 0xba, 0xd2, 0x01, 0x00, 0x00, 0x60, 0x01, 0x00})
			resp := &RawSendOutput{}
			result := ReadSerial(port)
			resp.Body.Message = string(result)
			return resp, nil
		})

		huma.Post(api, "/power/set/off", func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
			WriteSerial(port, []byte{0xbe, 0xef, 0x03, 0x06, 0x00, 0x2a, 0xd3, 0x01, 0x00, 0x00, 0x60, 0x00, 0x00})
			resp := &RawSendOutput{}
			result := ReadSerial(port)
			resp.Body.Message = string(result)
			return resp, nil
		})

		huma.Get(api, "/power/get", func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
			WriteSerial(port, []byte{0xbe, 0xef, 0x03, 0x06, 0x00, 0x19, 0xd3, 0x02, 0x00, 0x00, 0x60, 0x00, 0x00})
			resp := &RawSendOutput{}
			result := ReadSerial(port)
			resp.Body.Message = string(result)
			return resp, nil
		})

		huma.Post(api, "/input/set/rgb1", func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
			WriteSerial(port, []byte{0xbe, 0xef, 0x03, 0x06, 0x00, 0xfe, 0xd2, 0x01, 0x00, 0x00, 0x20, 0x00, 0x00})
			resp := &RawSendOutput{}
			result := ReadSerial(port)
			resp.Body.Message = string(result)
			return resp, nil
		})

		huma.Post(api, "/input/set/rgb2", func(ctx context.Context, input *struct{}) (*RawSendOutput, error) {
			WriteSerial(port, []byte{0xbe, 0xef, 0x03, 0x06, 0x00, 0x3e, 0xd0, 0x01, 0x00, 0x00, 0x20, 0x04, 0x00})
			resp := &RawSendOutput{}
			result := ReadSerial(port)
			resp.Body.Message = string(result)
			return resp, nil
		})

		fmt.Println("Routes registered")

		// Create the HTTP server.
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", options.Port),
			Handler: router,
		}

		// Tell the CLI how to start your server.
		hooks.OnStart(func() {
			fmt.Printf("Starting server on port %d...\n", options.Port)
			server.ListenAndServe()
		})

		// Tell the CLI how to stop your server.
		hooks.OnStop(func() {
			// Give the server 5 seconds to gracefully shut down, then give up.
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			fmt.Println("Disconnecting from serial...")
			DisconnectSerial(port)
			fmt.Println("Closed server successfully")
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
