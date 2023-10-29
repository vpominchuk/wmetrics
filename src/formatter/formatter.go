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
	fmt.Printf(strPadRight("Total time taken for tests:", strLength)+"%s\n", toTimeString(stat.TotalTime))
	fmt.Printf(strPadRight("Time per request (avg):", strLength)+"%.3f ms\n", toMilliseconds(stat.TotalTimeAvg))
	fmt.Printf(strPadRight("Time per request (median):", strLength)+"%.3f ms\n", toMilliseconds(stat.TotalTimeMedian))
	fmt.Printf(strPadRight("Time per request (min):", strLength)+"%.3f ms\n", toMilliseconds(stat.TotalTimeMin))
	fmt.Printf(strPadRight("Time per request (max):", strLength)+"%.3f ms\n", toMilliseconds(stat.TotalTimeMax))
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
		strPadRight("\nConnection Metrics:", strLength+3) +
			strPadRight("(avg)", 13) +
			strPadRight("(median)", 17) +
			strPadRight("(min)", 15) +
			strPadRight("(max)", 15),
	)

	printDurations(
		"DNS lookup:",
		toTimeString(stat.DNSLookupAvg),
		toTimeString(stat.DNSLookupMedian),
		toTimeString(stat.DNSLookupMin),
		toTimeString(stat.DNSLookupMax),
		strLength,
	)

	printDurations(
		"TCP connection:",
		toTimeString(stat.TCPConnectionAvg),
		toTimeString(stat.TCPConnectionMedian),
		toTimeString(stat.TCPConnectionMin),
		toTimeString(stat.TCPConnectionMax),
		strLength,
	)

	printDurations(
		"TLS handshake:",
		toTimeString(stat.TLSHandshakeAvg),
		toTimeString(stat.TLSHandshakeMedian),
		toTimeString(stat.TLSHandshakeMin),
		toTimeString(stat.TLSHandshakeMax),
		strLength,
	)

	printDurations(
		"Connection established:",
		toTimeString(stat.ConnectionEstablishedAvg),
		toTimeString(stat.ConnectionEstablishedMedian),
		toTimeString(stat.ConnectionEstablishedMin),
		toTimeString(stat.ConnectionEstablishedMax),
		strLength,
	)

	printDurations(
		"TTFB:",
		toTimeString(stat.TTFBAvg),
		toTimeString(stat.TTFBMedian),
		toTimeString(stat.TTFBMin),
		toTimeString(stat.TTFBMax),
		strLength,
	)

	if stat.TotalTimePercentage != nil && len(stat.TotalTimePercentage) > 0 {
		fmt.Println("\nPercentage of the requests served within a certain time (ms):")

		for _, result := range stat.TotalTimePercentage {
			fmt.Print(strPadRight(fmt.Sprintf("%d%%", result.Segment*10), 7))
			fmt.Print(strPadRight(fmt.Sprintf("%.3f ms", toMilliseconds((result.Min+result.Max)/2)), 15))
			fmt.Printf("(%.3f - %.3f ms)\n", toMilliseconds(result.Min), toMilliseconds(result.Max))
		}
	}

	if stat.Errors != nil && len(stat.Errors) > 0 {
		fmt.Println("\nErrors:")

		for _, result := range stat.Errors {
			fmt.Printf("%s (%d times)\n", result.Message, result.Count)
		}
	}
}

func toMilliseconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Millisecond)
}

func toSeconds(duration time.Duration) float64 {
	return float64(duration) / float64(time.Second)
}

func toTimeString(duration time.Duration) string {
	if duration < time.Second {
		return fmt.Sprintf("%.3f ms", toMilliseconds(duration))
	} else {
		return fmt.Sprintf("%.3f s", toSeconds(duration))
	}
}

func strPadRight(string string, count int) string {
	padLength := count - len(string)

	if padLength <= 0 {
		padLength = 1
	}

	return string + strings.Repeat(" ", padLength)
}

func printDurations(title string, avg, median, min, max string, strLength int) {
	fmt.Println(
		strPadRight(title, strLength) +
			strPadRight(fmt.Sprintf("%s", avg), 15) +
			strPadRight(fmt.Sprintf("%s", median), 15) +
			strPadRight(fmt.Sprintf("%s", min), 15) +
			strPadRight(fmt.Sprintf("%s", max), 15),
	)
}
