package formatter

import (
	"fmt"
	"strings"
	"time"
	"webmetrics/wmetrics/src/statistics"
)

func PrintResults(stat statistics.Statistics) {
	strLength := 30

	fmt.Printf(strPadRight("Complete requests:", strLength)+"%d\n", stat.TotalRequests)
	fmt.Printf(strPadRight("Successful requests:", strLength)+"%d\n", stat.SuccessRequests)
	fmt.Printf(strPadRight("Failed requests:", strLength)+"%d\n", stat.ErrorRequests)

	fmt.Println("\nPerformance Metrics:")
	fmt.Printf(strPadRight("Total time taken for tests:", strLength)+"%.3f seconds\n", toSeconds(stat.TotalTime))
	fmt.Printf(strPadRight("Time per request (avg):", strLength)+"%.3f ms\n", toMilliseconds(stat.TotalTimeAvg))
	fmt.Printf(strPadRight("Min time per request:", strLength)+"%.3f ms\n", toMilliseconds(stat.TotalTimeMin))
	fmt.Printf(strPadRight("Max time per request:", strLength)+"%.3f ms\n", toMilliseconds(stat.TotalTimeMax))
	fmt.Printf(
		strPadRight("Requests per second:", strLength)+"%.2f\n", float64(stat.TotalRequests)/toSeconds(stat.TotalTime),
	)

	if stat.Code2xx > 0 {
		fmt.Printf(strPadRight("2xx responses:", strLength)+"%d\n", stat.Code2xx)
	}

	if stat.Code3xx > 0 {
		fmt.Printf(strPadRight("3xx responses:", strLength)+"%d\n", stat.Code3xx)
	}

	if stat.Code5xx > 0 {
		fmt.Printf(strPadRight("5xx responses:", strLength)+"%d\n", stat.Code5xx)
	}

	if stat.OtherCodes > 0 {
		fmt.Printf(strPadRight("Other responses:", strLength)+"%d\n", stat.OtherCodes)
	}

	fmt.Println(
		strPadRight("\nConnection Metrics:", strLength+3) + strPadRight("(avg)", 12) +
			strPadRight("(min)", 12) +
			strPadRight("(max)", 12),
	)

	printDurations(
		"DNS lookup time:",
		toMilliseconds(stat.DNSLookupAvg),
		toMilliseconds(stat.DNSLookupMin),
		toMilliseconds(stat.DNSLookupMax),
		strLength,
	)
	// fmt.Printf(strPadRight("DNS lookup time:", strLength))
	// fmt.Print(strPadRight(fmt.Sprintf("%.3f ms", toMilliseconds(stat.DNSLookupAvg)), 12))
	// fmt.Print(strPadRight(fmt.Sprintf("%.3f ms", toMilliseconds(stat.DNSLookupMin)), 12))
	// fmt.Print(strPadRight(fmt.Sprintf("%.3f ms\n", toMilliseconds(stat.DNSLookupMax)), 12))
}

func toMilliseconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Millisecond)
}

func toSeconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Second)
}

func strPadRight(string string, count int) string {
	padLength := count - len(string)

	if padLength <= 0 {
		padLength = 1
	}

	return string + strings.Repeat(" ", padLength)
}

func printDurations(title string, avg, min, max float64, strLength int) {
	fmt.Println(
		strPadRight("DNS lookup time:", strLength) +
			strPadRight(fmt.Sprintf("%.3f ms", avg), 12) +
			strPadRight(fmt.Sprintf("%.3f ms", min), 12) +
			strPadRight(fmt.Sprintf("%.3f ms\n", max), 12),
	)
}
