package main

import (
	"fmt"
	"os"
	commandLine "webmetrics/wmetrics/src/args"
	"webmetrics/wmetrics/src/tester"
)

func main() {
	arguments, args := commandLine.GetArguments()

	if len(args) == 0 || len(args) > 1 {
		commandLine.Usage()
		os.Exit(1)
	}

	fmt.Printf("Flags: %v\n", *arguments.Method.Value)
	fmt.Printf("Args: %v\n", args)

	parameters := tester.Parameters{
		Resource:    args[0],
		Requests:    *arguments.Requests.Value,
		Concurrency: *arguments.Concurrency.Value,
		Timeout:     *arguments.Timeout.Value,
		Method:      *arguments.Method.Value,
		UserAgent:   *arguments.UserAgent.Value,
	}

	err := tester.Test(parameters)

	if err != nil {
		panic(err)
	}

	// https://pkg.go.dev/net/http#Transport
}
