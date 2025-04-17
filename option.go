package spin

import "github.com/FollowTheProcess/hue"

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
