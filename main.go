package main

import (
	"fmt"
	"net/url"
	"os"
	commandLine "webmetrics/wmetrics/src/args"
	"webmetrics/wmetrics/src/formatter"
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

	parsedUrl, err := url.Parse(args[0])

	if err != nil {
		panic(err)
	}

	parameters := tester.Parameters{
		Resource:              args[0],
		Url:                   parsedUrl,
		Requests:              *arguments.Requests.Value,
		Concurrency:           *arguments.Concurrency.Value,
		Timeout:               *arguments.Timeout.Value,
		Method:                *arguments.Method.Value,
		UserAgent:             *arguments.UserAgent.Value,
		UserAgentTemplate:     *arguments.UserAgentTemplate.Value,
		KeepAlive:             *arguments.KeepAlive.Value,
		Proxy:                 *arguments.Proxy.Value,
		MaxIdleConnections:    *arguments.MaxIdleConnections.Value,
		IdleConnTimeout:       *arguments.IdleConnTimeout.Value,
		TLSHandshakeTimeout:   *arguments.TLSHandshakeTimeout.Value,
		IPv4Only:              *arguments.IPv4Only.Value,
		IPv6Only:              *arguments.IPv6Only.Value,
		AllowInsecureSSL:      *arguments.AllowInsecureSSL.Value,
		ClientCertificateFile: *arguments.ClientCertificateFile.Value,
	}

	result, err := tester.Test(parameters)

	if err != nil {
		panic(err)
	}

	formatter.PrintResults(result)

	// https://pkg.go.dev/net/http#Transport
}
