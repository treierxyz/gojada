# gojada
> REST API for serial devices, built with [Go](https://go.dev/) and [Huma](https://huma.rocks/)

### Supported devices:
| Device | Identifier | Description                     |
|--------|------------|---------------------------------|
| 3M X55 | `x55`      | Projector with serial interface |

## Setup
1. In the repository root, run `go mod tidy` to install dependencies.
2. Start the server with `go run .`

### CLI options
| Option     | Description           | Default value  |
|------------|-----------------------|----------------|
| `--port`   | Port to listen on     | `8888`         |
| `--path`   | Path of serial device | `/dev/ttyUSB0` |
| `--device` | Device type           | `x55`          |

You can run the server with 'hot reload' using `air`:

0. `go install github.com/air-verse/air@latest`
1. `air`