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
	"github.com/treierxyz/hk-projector-api/devices/x55"
	"github.com/treierxyz/hk-projector-api/route"
	"github.com/treierxyz/hk-projector-api/serialhelper"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Options for the CLI.
type Options struct {
	Port   int    `help:"Port to listen on" short:"p" default:"8888"`
	Serial string `help:"Path of emulated device serial" short:"E" default:"/dev/ttyUSB0"`
}

func main() {
	// Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		fmt.Println("Connecting to serial...")
		port := serialhelper.ConnectSerial(options.Serial, x55.Device.Settings())
		fmt.Println("Connected successfully to serial ports")

		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("Häkkerikoda Projector API", "0.2.0"))

		route.CreateRoutes(api, port, x55.Device)

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
			// Give the server 15 seconds to gracefully shut down, then give up.
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			server.Shutdown(ctx)
			fmt.Println("Disconnecting from serial...")
			serialhelper.DisconnectSerial(port)
			fmt.Println("Closed server successfully")
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
