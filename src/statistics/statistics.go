package statistics

import (
	"time"
	"webmetrics/wmetrics/src/tester"
)

type Statistics struct {
	TotalTimeAvg,
	TotalTimeMin,
	TotalTimeMax,
	TotalTime,
	DNSLookupAvg,
	DNSLookupMin,
	DNSLookupMax time.Duration

	ErrorRequests,
	SuccessRequests,
	TotalRequests,
	Code2xx,
	Code3xx,
	Code5xx,
	OtherCodes int
}

func GetStatistics(results []tester.MeasurementResult, testDuration time.Duration) (Statistics, error) {
	var requestTimeAvg, requestTimeMin, requestTimeMax time.Duration
	var dnsLookupAvg, dnsLookupMin, dnsLookupMax time.Duration
	var errorRequests, successRequests, totalRequests, code2xx, code3xx, code5xx, otherCodes int

	for _, result := range results {
		if result.Error == nil {
			requestTimeAvg += result.RequestResult.Durations.Total.Total
			requestTimeMax = maxDuration(requestTimeMax, result.RequestResult.Durations.Total.Total)
			requestTimeMin = minDuration(requestTimeMin, result.RequestResult.Durations.Total.Total)

			dnsLookupAvg += result.RequestResult.Durations.DNSLookup.Duration
			dnsLookupMax = maxDuration(dnsLookupMax, result.RequestResult.Durations.DNSLookup.Duration)
			dnsLookupMin = minDuration(dnsLookupMin, result.RequestResult.Durations.DNSLookup.Duration)

			successRequests++

			if result.RequestResult.StatusCode >= 200 && result.RequestResult.StatusCode < 300 {
				code2xx++
			} else if result.RequestResult.StatusCode >= 300 && result.RequestResult.StatusCode < 400 {
				code3xx++
			} else if result.RequestResult.StatusCode >= 500 && result.RequestResult.StatusCode < 600 {
				code5xx++
			} else {
				otherCodes++
			}
		} else {
			errorRequests++
		}

		totalRequests++
	}

	return Statistics{
		TotalTimeAvg: requestTimeAvg / time.Duration(len(results)),
		TotalTimeMin: requestTimeMin,
		TotalTimeMax: requestTimeMax,

		DNSLookupAvg: dnsLookupAvg / time.Duration(len(results)),
		DNSLookupMax: dnsLookupMax,
		DNSLookupMin: dnsLookupMin,

		TotalTime:       testDuration,
		ErrorRequests:   errorRequests,
		SuccessRequests: successRequests,
		TotalRequests:   totalRequests,
		Code2xx:         code2xx,
		Code3xx:         code3xx,
		Code5xx:         code5xx,
		OtherCodes:      otherCodes,
	}, nil
}

func maxDuration(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}

	return b
}

func minDuration(a, b time.Duration) time.Duration {
	if b < a || a == 0 {
		return b
	}

	return a
}
