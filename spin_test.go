package spin_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/FollowTheProcess/spin"
)

func TestSpinner(t *testing.T) {
	buf := &bytes.Buffer{}

	spinner := spin.New(buf, "Testing")
	spinner.Start()
	time.Sleep(300 * time.Millisecond)
	spinner.Stop()

	got := buf.String()

	frames := [...]rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

	found := false
	for _, frame := range frames {
		if strings.Contains(got, string(frame)) {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected output to contain a spinner frame:\n\n%s\n", got)
	}
}
