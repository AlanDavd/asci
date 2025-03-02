# ASCI - Image to ASCII Art Converter

A cross-platform image to ASCII art converter written in Go with a hexagonal architecture. The project provides three ways to use the converter:

1. As a CLI tool
2. As an HTTP server with a web interface
3. As a Go module that can be imported into other Go programs

## Features

- Convert images (JPEG, PNG) to ASCII art
- Customizable output width and height
- Configurable character set for ASCII art
- Optional colored output (ANSI colors)
- Brightness inversion option
- Output to console or file
- Web interface for easy testing

## Installation

```bash
go install github.com/alandavd/asci/cmd/cli@latest
go install github.com/alandavd/asci/cmd/server@latest
```

## Usage

### As a CLI Tool

```bash
# Basic usage
asci input.jpg

# Customize width
asci -width 120 input.jpg

# Save to file
asci -output result.txt input.jpg

# Enable colored output
asci -color input.jpg

# Invert brightness
asci -invert input.jpg

# Full options
asci -width 120 -height 60 -charset "@%#*+=-:. " -color -invert -output result.txt input.jpg
```

### As an HTTP Server

```bash
# Start the server
server

# Server will be available at http://localhost:8080
```

Visit http://localhost:8080 in your browser to use the web interface.

### As a Go Module

```go
package main

import (
    "os"
    "github.com/alandavd/asci/pkg"
)

func main() {
    // Create a converter
    converter := asci.NewConverter()

    // Open an image file
    file, _ := os.Open("input.jpg")
    defer file.Close()

    // Convert with options
    result, err := converter.Convert(file,
        asci.WithWidth(120),
        asci.WithHeight(60),
        asci.WithCharset(" .:-=+*#%@"),
        asci.WithColor(true),
        asci.WithInverted(false),
    )

    if err != nil {
        panic(err)
    }

    // Print the result
    println(result)
}
```

## API Reference

### Options

- `WithWidth(width int)`: Set the width of the ASCII art output
- `WithHeight(height int)`: Set the height of the ASCII art output (0 for auto)
- `WithCharset(charset string)`: Set the characters to use for ASCII art
- `WithColor(enabled bool)`: Enable or disable colored output
- `WithInverted(enabled bool)`: Enable or disable brightness inversion

## Architecture

The project follows a hexagonal architecture pattern:

- `pkg/`: Public API for external use
- `internal/core/ports/`: Core interfaces
- `internal/core/services/`: Business logic implementation
- `cmd/cli/`: CLI application
- `cmd/server/`: HTTP server application

## License

MIT License
