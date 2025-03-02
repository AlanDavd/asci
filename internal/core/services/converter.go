package services

import (
	"image"
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/alandavd/asci/internal/core/ports"
)

type asciiConverter struct {
	options *ports.ConvertOptions
}

// NewASCIIConverter creates a new instance of the ASCII converter
func NewASCIIConverter() ports.ImageToASCIIConverter {
	return &asciiConverter{
		options: ports.DefaultOptions(),
	}
}

// Convert implements ports.ImageToASCIIConverter
func (c *asciiConverter) Convert(input io.Reader, options ...ports.ConvertOption) (string, error) {
	// Apply options
	opts := *c.options
	for _, option := range options {
		option(&opts)
	}

	// Decode image
	img, _, err := image.Decode(input)
	if err != nil {
		return "", err
	}

	// Convert to ASCII
	return c.imageToASCII(img, &opts), nil
}

// ConvertToFile implements ports.ImageToASCIIConverter
func (c *asciiConverter) ConvertToFile(input io.Reader, outputPath string, options ...ports.ConvertOption) error {
	ascii, err := c.Convert(input, options...)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, []byte(ascii), 0644)
}

func (c *asciiConverter) imageToASCII(img image.Image, opts *ports.ConvertOptions) string {
	// Get image bounds
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	// Calculate aspect ratio and resize dimensions
	aspectRatio := float64(height) / float64(width)
	newWidth := opts.Width
	newHeight := int(float64(opts.Width) * aspectRatio * 0.5) // Multiply by 0.5 to account for terminal font aspect ratio

	if opts.Height > 0 {
		newHeight = opts.Height
	}

	// Convert to ASCII
	var builder strings.Builder
	charsetLen := len(opts.Charset)
	
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Map image coordinates to original image
			srcX := x * width / newWidth
			srcY := y * height / newHeight
			
			// Get pixel color
			c := img.At(srcX, srcY)
			r, g, b, _ := c.RGBA()
			
			// Calculate brightness (0-255)
			brightness := (r + g + b) / 3 / 256
			
			// Map brightness to character
			if opts.Inverted {
				brightness = 65535 - brightness
			}
			
			charIndex := int(float64(brightness) * float64(charsetLen-1) / 65535)
			char := string(opts.Charset[charIndex])
			
			if opts.Colored {
				// Add ANSI color codes here if colored output is desired
				builder.WriteString("\x1b[38;2;" + strconv.Itoa(int(r/256)) + ";" + strconv.Itoa(int(g/256)) + ";" + strconv.Itoa(int(b/256)) + "m")
			}
			
			builder.WriteString(char)
		}
		builder.WriteString("\n")
	}
	
	if opts.Colored {
		builder.WriteString("\x1b[0m") // Reset colors
	}
	
	return builder.String()
} 