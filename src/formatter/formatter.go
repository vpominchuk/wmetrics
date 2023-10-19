package formatter

import (
	"fmt"
	"webmetrics/wmetrics/src/tester"
)

func PrintResults(result tester.MeasurementsResult) {
	if result.TLS.UseTLS {
		println(result.TLS.TLSVersion)
	}

	fmt.Printf("%v", result)
}
