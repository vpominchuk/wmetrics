package formatter

import (
	"fmt"
	"webmetrics/wmetrics/src/tester"
)

func PrintResults(result tester.MeasurementsResult) {
	resultJson, _ := result.ToJson()

	fmt.Printf("%s", resultJson)
}
