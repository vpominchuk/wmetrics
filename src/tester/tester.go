package tester

import (
	"errors"
	"time"
)

var testers = map[string]TestEngine{
	"http":  &HttpEngine{},
	"https": &HttpEngine{},
}

func Test(parameters Parameters, onProgress func(progress RequestsProgress)) (
	[]MeasurementResult, time.Duration, error,
) {
	testService, ok := testers[parameters.Url.Scheme]

	if ok {
		measurementResult, duration := testService.Measure(parameters, onProgress)
		return measurementResult, duration, nil
	}

	return []MeasurementResult{}, 0, errors.New("unsupported protocol")
}
