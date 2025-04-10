package spin_test

import (
	"testing"

	"github.com/FollowTheProcess/spin"
)

func TestHello(t *testing.T) {
	got := spin.Hello()
	want := "Hello spin"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}
