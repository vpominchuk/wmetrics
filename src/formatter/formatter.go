package formatter

import (
	"fmt"
	"time"
	"webmetrics/wmetrics/src/tester"
)

func PrintResults(result tester.MeasurementsResult) {
	resultJson, _ := result.ToJson()

	fmt.Printf("%s\n", resultJson)

	fmt.Printf("Total time: %.2f ms\n", toMilliseconds(result.Durations.Total.Total))
}

func toMilliseconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Millisecond)
}
