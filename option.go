package spin

import "go.followtheprocess.codes/hue"

// Option is a functional option for configuring a [Spinner].
type Option func(s *Spinner)

// MessageStyle sets the colour/style of the message text.
func MessageStyle(style hue.Style) Option {
	return func(s *Spinner) {
		s.messageStyle = style
	}
}

// FrameStyle sets the colour/style of the spinner animation frames.
func FrameStyle(style hue.Style) Option {
	return func(s *Spinner) {
		s.frameStyle = style
	}
}

// WithForceEnabled bypasses the terminal check, causing the spinner to render
// even when the writer is not a TTY. Useful for CI environments or testing.
func WithForceEnabled() Option {
	return func(s *Spinner) {
		s.forceEnabled = true
	}
}
