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

	results, testDuration, err := tester.Test(
		parameters,
		func(progress tester.RequestsProgress) {
			if bar != nil {
				bar.Set(progress.CompletedRequests)
			}
		},
	)

	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if len(results) == 0 {
		log.Fatalf("Error: something went wrong. No test results\n")
	}

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

	if len(args) == 0 {
		commandLine.Usage()
		os.Exit(1)
	}

	resources := make([]tester.Resource, 0, len(args))

	for _, link := range args {
		parsedUrl, err := url.Parse(link)

		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}

		resources = append(
			resources, tester.Resource{
				Url: parsedUrl,
			},
		)
	}

	return tester.Parameters{
		Resources:             resources,
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
		CustomHeaders:         *arguments.CustomHeaders.Value,
		TimeLimit:             *arguments.TimeLimit.Value,
	}
}

func buildProgressBar(parameters tester.Parameters) *progressbar.ProgressBar {
	progressBarMax := parameters.Requests

	if parameters.TimeLimit > 0 {
		progressBarMax = -1
	}

	return progressbar.NewOptions(
		progressBarMax,
		progressbar.OptionFullWidth(),
		progressbar.OptionShowCount(),
	)
}

func showGreetings(parameters tester.Parameters) {
	fmt.Printf("%s %s\n", app.ExecutableName, app.VersionString)
	fmt.Printf("Copyright %d Vasyl Pominchuk\n", time.Now().Year())

	if parameters.TimeLimit > 0 {
		fmt.Printf(
			"Performing requests with concurrency level of %d with time limit of %s\n",
			parameters.Concurrency,
			parameters.TimeLimit,
		)
	} else {
		fmt.Printf(
			"Performing %d requests with concurrency level of %d\n",
			parameters.Requests,
			parameters.Concurrency,
		)
	}

	fmt.Printf("%s", parameters.Method)

	if parameters.Resources != nil && len(parameters.Resources) > 1 {
		fmt.Printf(" %s\n", parameters.Resources[0].Url.String())
	} else {
		fmt.Printf(" List of URLs\n")
	}

	fmt.Printf("\n")
}
