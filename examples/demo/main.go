package main

import (
	"os"
	"time"

	"go.followtheprocess.codes/spin"
)

const duration = 2 * time.Second

func main() {
	spinner := spin.New(os.Stdout, "Digesting")

	spinner.Start()
	defer spinner.Stop()

	time.Sleep(duration)
}
