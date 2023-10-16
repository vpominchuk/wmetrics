package tester

import (
	"net/url"
)

var testers = map[string]TestEngine{
	"http":  &HttpEngine{},
	"https": &HttpEngine{},
}

func Test(parameters Parameters) error {
	resource, err := url.Parse(parameters.Resource)

	if err != nil {
		return err
	}

	testService, ok := testers[resource.Scheme]

	if ok {
		testService.Test(parameters)
	}

	return nil
}
