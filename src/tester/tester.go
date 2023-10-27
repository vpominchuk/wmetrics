package tester

import (
	"net/url"
)

var testers = map[string]TestEngine{
	"http":  &HttpEngine{},
	"https": &HttpEngine{},
}

func Test(parameters Parameters) []MeasurementResult {
	resource, err := url.Parse(parameters.Resource)

	if err != nil {
		return []MeasurementResult{}
	}

	testService, ok := testers[resource.Scheme]

	if ok {
		return testService.Measure(parameters)
	}

	return []MeasurementResult{}
}
