package main

import (
	"bufio"
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

	correctNumberOfRequests(&parameters)

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
				err := bar.Set(progress.CompletedRequests)

				if err != nil {
					return
				}
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

	if haveErrors(stat) {
		os.Exit(1)
	}

	os.Exit(0)
}

func haveErrors(stat statistics.Statistics) bool {
	for _, singleUrlStat := range stat {
		if len(singleUrlStat.Errors) > 0 {
			return true
		}
	}

	return false
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
	arguments, urls := commandLine.GetArguments()

	if err := commandLine.Validate(arguments); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if arguments.URLListFile.Value != nil && *arguments.URLListFile.Value != "" {
		var err error
		urls, err = getUrlsFromFile(*arguments.URLListFile.Value)

		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	}

	if len(urls) == 0 {
		commandLine.Usage()
		os.Exit(1)
	}

	resources := make([]tester.Resource, 0, len(urls))

	for _, link := range urls {
		parsedUrl, err := url.ParseRequestURI(link)

		if err != nil || parsedUrl.Scheme == "" || parsedUrl.Host == "" {
			stdError(fmt.Sprintf("* Warning: Skipping invalid url: %s\n", link))
			continue
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
		ExitWithErrorOnCode:   *arguments.ExitWithErrorOnCode.Value,
	}
}

func stdError(message string) {
	fmt.Fprint(os.Stderr, message)
}

func correctNumberOfRequests(parameters *tester.Parameters) {
	parameters.Requests = len(parameters.Resources) * parameters.Requests
}

func getUrlsFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	urls := make([]string, 0)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
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
			"Performing [%s] requests with concurrency level of %d with time limit of %s\n",
			parameters.Method,
			parameters.Concurrency,
			parameters.TimeLimit,
		)
	} else {
		fmt.Printf(
			"Performing %d [%s] requests with concurrency level of %d\n",
			parameters.Requests,
			parameters.Method,
			parameters.Concurrency,
		)
	}

	fmt.Printf("\n")
}
