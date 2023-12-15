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
	testService, ok := testers[parameters.Resources[0].Url.Scheme]

	resourceFeeder := newResourceFeeder(parameters.Resources)

	if ok {
		measurementResult, duration := testService.Measure(parameters, resourceFeeder, onProgress)
		return measurementResult, duration, nil
	}

	return []MeasurementResult{}, 0, errors.New("unsupported protocol")
}

func newResourceFeeder(resources []Resource) *ResourceFeeder {
	return &ResourceFeeder{
		Resources: resources,
		index:     0,
	}
}

func (s *ResourceFeeder) GetNextValue() (Resource, error) {
	if len(s.Resources) == 0 {
		return Resource{}, errors.New("url list is empty. No resources to test")
	}

	if len(s.Resources) == 1 {
		return s.Resources[0], nil
	}

	value := s.Resources[s.index]
	s.index++

	if s.index >= len(s.Resources) {
		s.index = 0
	}

	return value, nil
}
