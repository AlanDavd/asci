package ports

import "io"

// ImageToASCIIConverter defines the main interface for converting images to ASCII art
type ImageToASCIIConverter interface {
	// Convert converts an image to ASCII art
	Convert(input io.Reader, options ...ConvertOption) (string, error)
	// ConvertToFile converts an image to ASCII art and saves it to a file
	ConvertToFile(input io.Reader, outputPath string, options ...ConvertOption) error
}

// ConvertOption defines functional options for the converter
type ConvertOption func(*ConvertOptions)

// ConvertOptions holds all the configuration for the conversion
type ConvertOptions struct {
	Width     int     // Maximum width of the ASCII art
	Height    int     // Maximum height of the ASCII art
	Charset   string  // Characters to use for ASCII art
	Colored   bool    // Whether to include ANSI color codes
	Inverted  bool    // Whether to invert the brightness
}

// DefaultOptions returns the default conversion options
func DefaultOptions() *ConvertOptions {
	return &ConvertOptions{
		Width:    80,
		Height:   0, // Auto-calculate based on aspect ratio
		Charset:  "@%#*+=-:. ", // Reversed charset for better visibility on light backgrounds
		Colored:  false,
		Inverted: false,
	}
}
