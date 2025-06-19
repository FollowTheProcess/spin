// Package spin provides a simple, easy to use terminal spinner with configurable styles
// to provide progress information in command line applications.
package spin // import "go.followtheprocess.codes/spin"

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"go.followtheprocess.codes/hue"
	"golang.org/x/term"
)

const (
	// frameRate is the rate at which spinner frames are rendered.
	frameRate = 100 * time.Millisecond

	// defaultMessageStyle is the default style for the spinner message.
	defaultMessageStyle = hue.Bold

	// defaultFrameStyle is the default style for the spinner frames.
	defaultFrameStyle = hue.Cyan

	// erase is the ANSI code for erasing the current line and resetting the cursor
	// position back to the start of the line.
	//
	// See https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797#erase-functions.
	erase = "\r\x1b[K"
)

// Spinner contains the spinner state.
type Spinner struct {
	w            io.Writer      // Where to draw the spinner
	stop         chan struct{}  // Signal for the spinner to stop
	msg          string         // Message to display during spinning e.g. "Loading"
	wg           sync.WaitGroup // Manages the rendering goroutine
	running      atomic.Bool    // Whether the spinner is currently running
	messageStyle hue.Style      // Colour/style for the spinner message
	frameStyle   hue.Style      // Colour/style for the spinner animation
}

// New returns a new [Spinner].
func New(w io.Writer, msg string, options ...Option) *Spinner {
	spinner := &Spinner{
		w:            w,
		stop:         make(chan struct{}),
		msg:          msg,
		messageStyle: defaultMessageStyle,
		frameStyle:   defaultFrameStyle,
	}

	for _, option := range options {
		option(spinner)
	}

	return spinner
}

// Start starts the spinner animation.
func (s *Spinner) Start() {
	if s.running.Load() || !isTerminal(s.w) {
		// If it's already running, or if it's not hooked up to a terminal
		// there's nothing for us to do
		return
	}

	s.running.Store(true)

	s.wg.Add(1)
	go func() {
		// Store the frames and the index locally so no need for synchronisation
		frames := [...]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		current := 0
		defer s.wg.Done()
		for {
			select {
			case <-s.stop:
				return
			case <-time.Tick(frameRate):
				fmt.Fprintf(s.w, "%s%s %s...", erase, s.frameStyle.Text(frames[current]), s.messageStyle.Text(s.msg))
				current = (current + 1) % len(frames)
			}
		}
	}()
}

// Stop halts the spinner animation.
func (s *Spinner) Stop() {
	if !s.running.Load() {
		// If it's not already running, there's nothing to do
		return
	}

	s.stop <- struct{}{} // Signal the stop
	s.wg.Wait()          // Wait for the render goroutine to return
	s.running.Store(false)

	fmt.Fprint(s.w, erase) // Erase the line
}

// Do runs the given function behind a spinner, automatically starting
// and stopping the spinner.
func (s *Spinner) Do(fn func()) {
	s.Start()
	defer s.Stop()
	fn()
}

// isTerminal returns whether w points to an [*os.File] and whether
// that is itself a tty.
func isTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}

	if file == nil {
		return false
	}

	return term.IsTerminal(int(file.Fd()))
}
