package asci

import (
	"github.com/alandavd/asci/internal/core/ports"
	"github.com/alandavd/asci/internal/core/services"
)

// Converter is the public interface for converting images to ASCII art
type Converter interface {
	ports.ImageToASCIIConverter
}

// NewConverter creates a new instance of the ASCII art converter
func NewConverter() Converter {
	return services.NewASCIIConverter()
}

// WithWidth sets the width of the ASCII art output
func WithWidth(width int) ports.ConvertOption {
	return func(o *ports.ConvertOptions) {
		o.Width = width
	}
}

// WithHeight sets the height of the ASCII art output
func WithHeight(height int) ports.ConvertOption {
	return func(o *ports.ConvertOptions) {
		o.Height = height
	}
}

// WithCharset sets the character set to use for the ASCII art
func WithCharset(charset string) ports.ConvertOption {
	return func(o *ports.ConvertOptions) {
		o.Charset = charset
	}
}

// WithColor enables or disables colored output
func WithColor(enabled bool) ports.ConvertOption {
	return func(o *ports.ConvertOptions) {
		o.Colored = enabled
	}
}

// WithInverted enables or disables brightness inversion
func WithInverted(enabled bool) ports.ConvertOption {
	return func(o *ports.ConvertOptions) {
		o.Inverted = enabled
	}
} 