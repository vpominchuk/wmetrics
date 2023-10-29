package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
	"net/url"
	"os"
	"strings"
	"time"
	"webmetrics/wmetrics/src/app"
	commandLine "webmetrics/wmetrics/src/args"
	"webmetrics/wmetrics/src/formatter"
	"webmetrics/wmetrics/src/statistics"
	"webmetrics/wmetrics/src/tester"
)

func main() {
	parameters := getCLIParameters()

	showGreetings(parameters)

	bar := buildProgressBar(parameters)

	results, testDuration := tester.Test(
		parameters,
		func(progress tester.RequestsProgress) {
			bar.Set(progress.CompletedRequests)
		},
	)

	stat, _ := statistics.GetStatistics(results, testDuration)

	fmt.Print("\n\n\n")

	formatter.PrintResults(stat)

	// for indx, result := range results {
	// 	fmt.Printf("Result: %d ", indx)
	//
	// 	if result.Error == nil {
	// 		formatter.PrintResults(result.RequestResult)
	// 	} else {
	// 		fmt.Printf("Error: %v\n", result.Error)
	// 	}
	// }
	fmt.Printf("\n")
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
