package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	commandLine "webmetrics/wmetrics/src/args"
	"webmetrics/wmetrics/src/formatter"
	"webmetrics/wmetrics/src/tester"
)

func main() {
	arguments, args := commandLine.GetArguments()

	jsonString, _ := json.MarshalIndent(arguments, "", "  ")
	fmt.Printf("%s\n", jsonString)

	if err := commandLine.Validate(arguments); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

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
		Method:                strings.ToUpper(*arguments.Method.Value),
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
		PostDataFile:          *arguments.PostDataFile.Value,
		PostData:              *arguments.PostData.Value,
		ContentType:           *arguments.ContentType.Value,
		FormData:              *arguments.FormData.Value,
	}

	result, err := tester.Test(parameters)

	if err != nil {
		panic(err)
	}

	formatter.PrintResults(result)

	// https://pkg.go.dev/net/http#Transport
}
