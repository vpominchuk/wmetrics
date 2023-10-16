package tester

import "fmt"

type HttpEngine struct {
}

func (engine *HttpEngine) Test(parameters Parameters) {
	fmt.Printf("HttpEngine: %v\n", parameters)
}
