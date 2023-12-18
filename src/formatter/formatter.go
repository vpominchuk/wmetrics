package formatter

import (
	"encoding/json"
	"fmt"
	"github.com/vpominchuk/wmetrics/src/statistics"
	"log"
	"strings"
	"time"
)

func PrintJsonResults(stat statistics.Statistics, pretty bool) {
	var jsonData []byte
	var err error

	if pretty {
		jsonData, err = json.MarshalIndent(stat, "", "  ")
	} else {
		jsonData, err = json.Marshal(stat)
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	fmt.Println(string(jsonData))
}

func printTitle(title string) {
	fmt.Println(title + ":")
}

func PrintResults(stats statistics.Statistics) {
	urlNum := 0

	for url, stat := range stats {
		printTitle(url)
		printSingleUrlResults(stat)

		if urlNum < len(stats)-1 {
			fmt.Print("─────────────────────────────────────────────────────────────────────────────────────\n\n")
		}

		urlNum++
	}
}

func printSingleUrlResults(stat statistics.SingleUrlStatistics) {
	strLength := 30

	if stat.Server != "" {
		fmt.Printf(StrPadRight("Server:", strLength)+"%s\n", stat.Server)
	}

	if stat.PoweredBy != "" {
		fmt.Printf(StrPadRight("Powered By:", strLength)+"%s\n", stat.PoweredBy)
	}

	fmt.Printf(StrPadRight("Complete requests:", strLength)+"%d\n", stat.TotalRequests)
	fmt.Printf(StrPadRight("Successful requests:", strLength)+"%d\n", stat.SuccessRequests)
	fmt.Printf(StrPadRight("Failed requests:", strLength)+"%d\n", stat.ErrorRequests)

	fmt.Println("\nPerformance Metrics:")
	fmt.Printf(StrPadRight("Total time taken for tests:", strLength)+"%s\n", toTimeString(stat.TotalTime))
	fmt.Printf(StrPadRight("Time per request (avg):", strLength)+"%s\n", toTimeString(stat.RequestTimeAvg))
	fmt.Printf(StrPadRight("Time per request (median):", strLength)+"%s\n", toTimeString(stat.RequestTimeMedian))
	fmt.Printf(StrPadRight("Time per request (min):", strLength)+"%s\n", toTimeString(stat.RequestTimeMin))
	fmt.Printf(StrPadRight("Time per request (max):", strLength)+"%s\n", toTimeString(stat.RequestTimeMax))
	fmt.Printf(
		StrPadRight("Requests per second:", strLength)+"%.2f\n", float64(stat.TotalRequests)/toSeconds(stat.TotalTime),
	)

	if stat.Code2xx > 0 {
		fmt.Printf(StrPadRight("2xx responses:", strLength)+"%d\n", stat.Code2xx)
	}

	if stat.Code3xx > 0 {
		fmt.Printf(StrPadRight("3xx responses:", strLength)+"%d\n", stat.Code3xx)
	}

	if stat.Code4xx > 0 {
		fmt.Printf(StrPadRight("4xx responses:", strLength)+"%d\n", stat.Code4xx)
	}

	if stat.Code5xx > 0 {
		fmt.Printf(StrPadRight("5xx responses:", strLength)+"%d\n", stat.Code5xx)
	}

	if stat.OtherCodes > 0 {
		fmt.Printf(StrPadRight("Other responses:", strLength)+"%d\n", stat.OtherCodes)
	}

	fmt.Println(
		StrPadRight("\nConnection Metrics:", strLength+3) +
			StrPadRight("(avg)", 13) +
			StrPadRight("(median)", 17) +
			StrPadRight("(min)", 15) +
			StrPadRight("(max)", 15),
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
			fmt.Print(StrPadRight(fmt.Sprintf("%d%%", result.Segment*10), 7))
			fmt.Print(StrPadRight(fmt.Sprintf("%.3f ms", toMilliseconds((result.Min+result.Max)/2)), 15))
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
	return fmt.Sprintf("%.3f ms", toMilliseconds(duration))
}

func StrPadRight(string string, count int) string {
	padLength := count - len(string)

	if padLength <= 0 {
		padLength = 1
	}

	return string + strings.Repeat(" ", padLength)
}

func printDurations(title string, avg, median, min, max string, strLength int) {
	fmt.Println(
		StrPadRight(title, strLength) +
			StrPadRight(avg, 15) +
			StrPadRight(median, 15) +
			StrPadRight(min, 15) +
			StrPadRight(max, 15),
	)
}
