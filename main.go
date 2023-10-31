package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/vpominchuk/wmetrics/src/app"
	commandLine "github.com/vpominchuk/wmetrics/src/args"
	"github.com/vpominchuk/wmetrics/src/formatter"
	"github.com/vpominchuk/wmetrics/src/statistics"
	"github.com/vpominchuk/wmetrics/src/tester"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	parameters := getCLIParameters()

	if canPrintGreetings(parameters.OutputFormat) {
		showGreetings(parameters)
	}

	var bar *progressbar.ProgressBar

	if canPrintProgressBar(parameters.OutputFormat) {
		bar = buildProgressBar(parameters)
	}

	results, testDuration := tester.Test(
		parameters,
		func(progress tester.RequestsProgress) {
			if bar != nil {
				bar.Set(progress.CompletedRequests)
			}
		},
	)

	stat, _ := statistics.GetStatistics(results, testDuration)

	if canPrintGreetings(parameters.OutputFormat) {
		fmt.Print("\n\n\n")
	}

	printResults(parameters.OutputFormat, stat)

	if canPrintGreetings(parameters.OutputFormat) {
		fmt.Printf("\n")
	}
}

func canPrintProgressBar(format string) bool {
	return strings.ToLower(format) == "std"
}

func canPrintGreetings(format string) bool {
	return strings.ToLower(format) == "std" || strings.ToLower(format) == "text"
}

func printResults(format string, stat statistics.Statistics) {
	switch strings.ToLower(format) {
	case "std", "text":
		formatter.PrintResults(stat)
	case "json":
		formatter.PrintJsonResults(stat, false)
	case "json-pretty":
		formatter.PrintJsonResults(stat, true)
	}
}

func getCLIParameters() tester.Parameters {
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
		log.Fatalf("Error: %v\n", err)
	}

	return tester.Parameters{
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
		OutputFormat:          *arguments.OutputFormat.Value,
	}
}

func buildProgressBar(parameters tester.Parameters) *progressbar.ProgressBar {
	return progressbar.NewOptions(
		parameters.Requests,
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
	)
}

func showGreetings(parameters tester.Parameters) {
	fmt.Printf("%s %s\n", app.ExecutableName, app.VersionString)
	fmt.Printf("Copyright %d Vasyl Pominchuk\n", time.Now().Year())
	fmt.Printf(
		"Performing %d requests with concurrency level of %d\n",
		parameters.Requests,
		parameters.Concurrency,
	)
	fmt.Printf(
		"%s %s\n",
		parameters.Method,
		parameters.Resource,
	)

	fmt.Printf("\n")
}
