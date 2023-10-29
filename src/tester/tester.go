package tester

import (
	"net/url"
	"time"
)

var testers = map[string]TestEngine{
	"http":  &HttpEngine{},
	"https": &HttpEngine{},
}

func Test(parameters Parameters, onProgress func(progress RequestsProgress)) ([]MeasurementResult, time.Duration) {
	resource, err := url.Parse(parameters.Resource)

	if err != nil {
		return []MeasurementResult{}, 0
	}

	testService, ok := testers[resource.Scheme]

	if ok {
		return testService.Measure(parameters, onProgress)
	}

	return []MeasurementResult{}, 0
}
