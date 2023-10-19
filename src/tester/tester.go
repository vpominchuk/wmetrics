package tester

import (
	"net/url"
)

var testers = map[string]TestEngine{
	"http":  &HttpEngine{},
	"https": &HttpEngine{},
}

func Test(parameters Parameters) (MeasurementsResult, error) {
	resource, err := url.Parse(parameters.Resource)

	if err != nil {
		return MeasurementsResult{}, err
	}

	testService, ok := testers[resource.Scheme]

	if ok {
		return testService.Measure(parameters)
	}

	return MeasurementsResult{}, nil
}
