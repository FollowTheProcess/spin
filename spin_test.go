package spin_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/FollowTheProcess/hue"
	"github.com/FollowTheProcess/spin"
)

var frames = [...]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func TestSpinner(t *testing.T) {
	buf := &bytes.Buffer{}

	spinner := spin.New(buf, "Testing", spin.FrameStyle(hue.Yellow), spin.MessageStyle(hue.Blue))
	spinner.Start()
	time.Sleep(300 * time.Millisecond)
	spinner.Stop()

	got := buf.String()

	found := false
	for _, frame := range frames {
		if strings.Contains(got, frame) {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected output to contain a spinner frame:\n\n%s\n", got)
	}
}

func TestDo(t *testing.T) {
	buf := &bytes.Buffer{}

	spinner := spin.New(buf, "Testing", spin.FrameStyle(hue.Yellow), spin.MessageStyle(hue.Blue))
	spinner.Do(func() {
		time.Sleep(300 * time.Millisecond)
	})

	got := buf.String()

	found := false
	for _, frame := range frames {
		if strings.Contains(got, frame) {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected output to contain a spinner frame:\n\n%s\n", got)
	}
}

func TestDoubleStart(t *testing.T) {
	buf := &bytes.Buffer{}

	spinner := spin.New(buf, "Working")
	spinner.Start()
	spinner.Start() // This one shouldn't do anything

	time.Sleep(300 * time.Millisecond)
	spinner.Stop()

	got := buf.String()

	found := false
	for _, frame := range frames {
		if strings.Contains(got, frame) {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected output to contain a spinner frame:\n\n%s\n", got)
	}
}

func TestDoubleStop(t *testing.T) {
	buf := &bytes.Buffer{}

	spinner := spin.New(buf, "Working")
	spinner.Start()

	time.Sleep(300 * time.Millisecond)
	spinner.Stop()
	spinner.Stop() // This one shouldn't do anything

	got := buf.String()

	found := false
	for _, frame := range frames {
		if strings.Contains(got, frame) {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected output to contain a spinner frame:\n\n%s\n", got)
	}
}
