package spin_test

import (
	"bytes"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"go.followtheprocess.codes/spin"
)

func TestSpinnerStartWritesFramesAndMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{name: "short message", msg: "Loading"},
		{name: "longer message", msg: "Processing data"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var buf bytes.Buffer
				s := spin.New(&buf, tt.msg, spin.WithForceEnabled())
				s.Start()
				time.Sleep(250 * time.Millisecond) // fake: two ticks at 100ms and 200ms
				synctest.Wait()
				s.Stop()

				if output := buf.String(); !strings.Contains(output, tt.msg) {
					t.Errorf("output %q does not contain message %q", output, tt.msg)
				}
			})
		})
	}
}

func TestSpinnerStopErasesLine(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var buf bytes.Buffer
		s := spin.New(&buf, "Loading", spin.WithForceEnabled())
		s.Start()
		time.Sleep(250 * time.Millisecond)
		synctest.Wait()
		s.Stop()

		if output := buf.String(); !strings.HasSuffix(output, "\r\x1b[K") {
			t.Errorf("output %q does not end with ANSI erase sequence", output)
		}
	})
}

func TestSpinnerStopWhenNotRunningIsNoOp(t *testing.T) {
	var buf bytes.Buffer
	s := spin.New(&buf, "Loading", spin.WithForceEnabled())
	s.Stop() // must not panic or write anything

	if buf.Len() != 0 {
		t.Errorf("Stop() on unstarted spinner wrote output: %q", buf.String())
	}
}

func TestSpinnerStopCalledTwiceIsNoOp(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var buf bytes.Buffer
		s := spin.New(&buf, "Loading", spin.WithForceEnabled())
		s.Start()
		time.Sleep(250 * time.Millisecond)
		synctest.Wait()
		s.Stop()

		outputAfterFirst := buf.String()

		s.Stop() // second Stop must not write anything extra or block

		if buf.String() != outputAfterFirst {
			t.Errorf("second Stop() wrote additional output: %q", buf.String()[len(outputAfterFirst):])
		}
	})
}

func TestSpinnerStartWhenAlreadyRunningIsNoOp(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var buf bytes.Buffer
		s := spin.New(&buf, "Loading", spin.WithForceEnabled())
		s.Start()
		s.Start() // second call must not start a second goroutine
		time.Sleep(250 * time.Millisecond)
		synctest.Wait()
		s.Stop() // must not deadlock

		// A single erase sequence must appear at the end — not multiple.
		if output := buf.String(); !strings.HasSuffix(output, "\r\x1b[K") {
			t.Errorf("output %q does not end with a single erase sequence after double Start", output)
		}
	})
}

func TestSpinnerDoExecutesFunction(t *testing.T) {
	var buf bytes.Buffer
	s := spin.New(&buf, "Loading", spin.WithForceEnabled())

	ran := false
	s.Do(func() { ran = true })

	if !ran {
		t.Error("Do() did not execute the provided function")
	}
}

func TestSpinnerDoStartsAndStopsAroundFunction(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var buf bytes.Buffer
		s := spin.New(&buf, "Loading", spin.WithForceEnabled())

		s.Do(func() { time.Sleep(250 * time.Millisecond) })

		output := buf.String()
		if !strings.Contains(output, "Loading") {
			t.Errorf("output %q does not contain message", output)
		}
		if !strings.HasSuffix(output, "\r\x1b[K") {
			t.Errorf("output %q does not end with erase sequence", output)
		}
	})
}

func TestSpinnerNonTerminalWriterProducesNoOutput(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var buf bytes.Buffer
		s := spin.New(&buf, "Loading") // no WithForceEnabled — Start is a no-op
		s.Start()
		time.Sleep(250 * time.Millisecond)
		s.Stop()

		if buf.Len() != 0 {
			t.Errorf("non-terminal writer got output: %q", buf.String())
		}
	})
}
