package main

import (
	"github.com/0xJeti/shuffledns/pkg/runner"
	"github.com/projectdiscovery/gologger"
)

func main() {
	// Parse the command line flags and read config files
	options := runner.ParseOptions()

	runner, err := runner.New(options)
	if err != nil {
		gologger.Fatalf("Could not create runner: %s\n", err)
	}

	runner.RunEnumeration()
	runner.Close()
}
