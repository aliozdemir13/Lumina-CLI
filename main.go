package main

import (
	"github.com/aliozdemir13/Lumina/cmd"
	"github.com/aliozdemir13/Lumina/internal"
)

func main() {
	// Print logo
	internal.PrintHeader()

	// Cobra takes over everything from here
	cmd.Execute()
}
