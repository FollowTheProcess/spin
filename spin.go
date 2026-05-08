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
	messageStyle hue.Style      // Colour/style for the spinner message
	frameStyle   hue.Style      // Colour/style for the spinner animation
	running      atomic.Bool    // Whether the spinner is currently running
	forceEnabled bool           // Bypass terminal check (useful for testing and CI)
}

// New returns a new [Spinner].
func New(w io.Writer, msg string, options ...Option) *Spinner {
	spinner := &Spinner{
		w:            w,
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
	if !s.forceEnabled && !isTerminal(s.w) {
		// Not hooked up to a terminal and not force-enabled; nothing to do.
		return
	}

	if !s.running.CompareAndSwap(false, true) {
		// Already running; a concurrent or duplicate Start call is a no-op.
		return
	}

	s.stop = make(chan struct{})

	s.wg.Go(func() {
		// Store the frames and the index locally so no need for synchronisation.
		frames := [...]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		current := 0
		ticker := time.NewTicker(frameRate)
		defer ticker.Stop()
		for {
			select {
			case <-s.stop:
				return
			case <-ticker.C:
				fmt.Fprintf(s.w, "%s%s %s...", erase, s.frameStyle.Text(frames[current]), s.messageStyle.Text(s.msg))
				current = (current + 1) % len(frames)
			}
		}
	})
}

// Stop halts the spinner animation.
func (s *Spinner) Stop() {
	if !s.running.CompareAndSwap(true, false) {
		// Not running; a concurrent or duplicate Stop call is a no-op.
		return
	}

	close(s.stop) // Signal the goroutine to stop.
	s.wg.Wait()   // Wait for the render goroutine to return.

	fmt.Fprint(s.w, erase) // Erase the line.
}

// Do runs the given function behind a spinner, automatically starting
// and stopping the spinner.
func (s *Spinner) Do(fn func()) {
	s.Start()
	defer s.Stop()
	fn()
}

// isTerminal returns whether w points to an [*os.File] that is itself a tty.
func isTerminal(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(file.Fd()))
}
