package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alandavd/asci/internal/core/ports"
	"github.com/alandavd/asci/pkg"
)

func main() {
	// Parse command line flags
	width := flag.Int("width", 80, "Width of the ASCII art output")
	height := flag.Int("height", 0, "Height of the ASCII art output (0 for auto)")
	charset := flag.String("charset", " .:-=+*#%@", "Characters to use for ASCII art")
	output := flag.String("output", "", "Output file path (optional)")
	colored := flag.Bool("color", false, "Enable colored output")
	inverted := flag.Bool("invert", false, "Invert brightness")
	flag.Parse()

	// Check if input file is provided
	if flag.NArg() != 1 {
		fmt.Println("Usage: asci [options] <input-file>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Open input file
	inputPath := flag.Arg(0)
	input, err := os.Open(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer input.Close()

	// Create converter with options
	converter := asci.NewConverter()
	options := []ports.ConvertOption{
		asci.WithWidth(*width),
		asci.WithCharset(*charset),
		asci.WithColor(*colored),
		asci.WithInverted(*inverted),
	}

	if *height > 0 {
		options = append(options, asci.WithHeight(*height))
	}

	// Convert image
	if *output != "" {
		// Convert to file
		err = converter.ConvertToFile(input, *output, options...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting image: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Convert and print to stdout
		result, err := converter.Convert(input, options...)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting image: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(result)
	}
} 