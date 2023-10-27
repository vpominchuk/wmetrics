package args

import (
	"fmt"
	"slices"
	"strings"
)

func Validate(arguments Arguments) error {
	postDataSources := getPostDataSourcesCount(arguments)

	if postDataSources > 1 {
		return fmt.Errorf("post data can only be specified once")
	}

	if *arguments.Requests.Value < *arguments.Concurrency.Value {
		return fmt.Errorf("cannot use concurrency level greater than total number of requests")
	}

	method := strings.ToUpper(*arguments.Method.Value)

	allowedMethods := []string{"GET", "HEAD", "DELETE", "POST", "PUT", "PATCH"}
	if !slices.Contains(allowedMethods, method) {
		return fmt.Errorf("invalid method: %s. Allowed methods are: %v", method, allowedMethods)
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

	return nil
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
