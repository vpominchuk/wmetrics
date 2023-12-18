package args

import (
	"fmt"
	"github.com/vpominchuk/wmetrics/src/app"
	"github.com/vpominchuk/wmetrics/src/formatter"
	"os"
	"regexp"
	"slices"
	"strings"
)

func Validate(arguments Arguments) error {
	postDataSources := getPostDataSourcesCount(arguments)

	if postDataSources > 1 {
		return fmt.Errorf("post data can only be specified once")
	}

	if *arguments.TimeLimit.Value == 0 && *arguments.Requests.Value < *arguments.Concurrency.Value {
		return fmt.Errorf("cannot use concurrency level greater than total number of requests")
	}

	method := strings.ToUpper(*arguments.Method.Value)

	allowedMethods := []string{"GET", "HEAD", "DELETE", "POST", "PUT", "PATCH"}
	if !slices.Contains(allowedMethods, method) {
		return fmt.Errorf("invalid method: %s. Allowed methods are: %v", method, allowedMethods)
	}

	allowedOutputFormats := []string{"std", "text", "json", "json-pretty"}
	if !slices.Contains(allowedOutputFormats, *arguments.OutputFormat.Value) {
		return fmt.Errorf(
			"invalid output format: %s. Allowed formats are: %v", *arguments.OutputFormat.Value, allowedOutputFormats,
		)
	}

	if arguments.FormData.Value != nil && *arguments.FormData.Value != "" {
		if arguments.ContentType.Value == nil || *arguments.ContentType.Value == "" {
			return fmt.Errorf("content type is required for form data. Use -%s", arguments.ContentType.Name)
		}

		if method != "POST" && method != "PUT" && method != "PATCH" {
			return fmt.Errorf(
				"method must be either POST, PUT or PATCH for form data. Current method: %s. Specify HTTP method using -%s METHOD",
				method, arguments.Method.Name,
			)
		}

		if *arguments.PostData.Value != "" {
			return fmt.Errorf("form data and post data cannot be used together")
		}

		if *arguments.PostDataFile.Value != "" {
			return fmt.Errorf("form data and post data file cannot be used together")
		}
	}

	if arguments.UserAgentTemplate.Value != nil && *arguments.UserAgentTemplate.Value != "" {
		if *arguments.UserAgentTemplate.Value == "list" {
			return fmt.Errorf("Allowed templates are: \n%s", getAllowedUserAgentTemplates())
		}

		if _, ok := app.DefaultUserAgents[*arguments.UserAgentTemplate.Value]; !ok {
			errorStr := fmt.Sprintf("invalid user agent template: %s.", *arguments.UserAgentTemplate.Value)

			return fmt.Errorf(errorStr)
		}
	}

	if err := validateUrlListFile(*arguments.URLListFile.Value); err != nil {
		return err
	}

	if arguments.ExitWithErrorOnCode.Value != nil && len(*arguments.ExitWithErrorOnCode.Value) > 0 {
		pattern := "(?i)^([0-9]{1,3}|[0-9]{1}xx)$"
		re, err := regexp.Compile(pattern)

		if err != nil {
			return fmt.Errorf("invalid exit with error on code pattern: %s", pattern)
		}

		for _, code := range *arguments.ExitWithErrorOnCode.Value {
			if !re.MatchString(code) {
				return fmt.Errorf("invalid exit with error on code: %s", code)
			}
		}
	}

	return nil
}

func validateUrlListFile(urlListFile string) error {
	if urlListFile == "" {
		return nil
	}

	if !fileExists(urlListFile) {
		return fmt.Errorf("url list file not found: %s", urlListFile)
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getAllowedUserAgentTemplates() string {
	result := ""

	for template, ua := range app.DefaultUserAgents {
		templateName := formatter.StrPadRight(template, 20)
		result += templateName + fmt.Sprintf(" %s\n", ua)
	}

	return result
}

func getPostDataSourcesCount(arguments Arguments) int {
	postDataSources := 0

	if *arguments.FormData.Value != "" {
		postDataSources++
	}

	if *arguments.PostData.Value != "" {
		postDataSources++
	}

	if *arguments.PostDataFile.Value != "" {
		postDataSources++
	}

	return postDataSources
}
