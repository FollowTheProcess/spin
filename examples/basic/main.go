package main

import (
	"os"
	"time"

	"github.com/FollowTheProcess/spin"
)

const duration = 2 * time.Second

func main() {
	spinner := spin.New(os.Stdout, "Digesting")

	spinner.Start()
	defer spinner.Stop()

	time.Sleep(duration)
}
