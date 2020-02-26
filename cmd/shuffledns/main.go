package main

import (
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/shuffledns/pkg/runner"
)

func main() {
	// Parse the command line flags and read config files
	options := runner.ParseOptions()

	runner, err := runner.NewRunner(options)
	if err != nil {
		gologger.Fatalf("Could not create runner: %s\n", err)
	}

}
