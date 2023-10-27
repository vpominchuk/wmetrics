package main

import (
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

	if err := commandLine.Validate(arguments); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if len(args) == 0 || len(args) > 1 {
		commandLine.Usage()
		os.Exit(1)
	}

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

	results := tester.Test(parameters)

	for indx, result := range results {
		fmt.Printf("Result: %d ", indx)

		formatter.PrintResults(result.RequestResult)
	}

	// if results[0].Error != nil {
	// 	panic(results[0].Error)
	// }
	//
	// formatter.PrintResults(results[0].RequestResult)

	// https://pkg.go.dev/net/http#Transport
}
