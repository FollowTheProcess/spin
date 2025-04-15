// Package spin provides a simple, easy to use terminal spinner with configurable styles
// to provide progress information in command line applications.
package spin

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// frameRate is the rate at which spinner frames are rendered.
	frameRate = 100 * time.Millisecond

	// erase is the ANSI code for erasing the current line and resetting the cursor
	// position back to the start of the line.
	//
	// See https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797#erase-functions.
	erase = "\r\x1b[K"
)

// Spinner contains the spinner state.
type Spinner struct {
	w       io.Writer      // Where to draw the spinner
	stop    chan struct{}  // Signal for the spinner to stop
	msg     string         // Message to display during spinning e.g. "Loading"
	wg      sync.WaitGroup // Manages the rendering goroutine
	running atomic.Bool    // Whether the spinner is currently running
}

// New returns a new [Spinner].
func New(w io.Writer, msg string) *Spinner {
	return &Spinner{
		w:    w,
		stop: make(chan struct{}),
		msg:  msg,
	}
}

// Start starts the spinner animation.
func (s *Spinner) Start() {
	if s.running.Load() {
		// If it's already running, we have nothing to do
		return
	}

	s.running.Store(true)

	s.wg.Add(1)
	go func() {
		// Store the frames and the index locally so no need for synchronisation
		frames := [...]rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
		current := 0
		defer s.wg.Done()
		for {
			select {
			case <-s.stop:
				return
			case <-time.Tick(frameRate):
				fmt.Fprintf(s.w, "%s%c %s...", erase, frames[current], s.msg)
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
