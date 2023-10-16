package tester

import "time"

type Parameters struct {
	Resource    string
	Requests    int
	Concurrency int
	Timeout     time.Duration
	Method      string
	UserAgent   string
}

type TestEngine interface {
	Test(parameters Parameters)
}
