package statistics

import (
	"github.com/vpominchuk/wmetrics/src/tester"
	"slices"
	"time"
)

type QuantileResult struct {
	Segment int
	Min     time.Duration
	Max     time.Duration
}

type ErrorResult struct {
	Message string
	Count   int
}

type SingleUrlStatistics struct {
	RequestTimeAvg,
	RequestTimeMin,
	RequestTimeMax,
	RequestTimeMedian,

	TotalTime,

	DNSLookupAvg,
	DNSLookupMin,
	DNSLookupMax,
	DNSLookupMedian,

	TCPConnectionAvg,
	TCPConnectionMin,
	TCPConnectionMax,
	TCPConnectionMedian,

	TLSHandshakeAvg,
	TLSHandshakeMin,
	TLSHandshakeMax,
	TLSHandshakeMedian,

	ConnectionEstablishedAvg,
	ConnectionEstablishedMin,
	ConnectionEstablishedMax,
	ConnectionEstablishedMedian,

	TTFBAvg,
	TTFBMin,
	TTFBMax,
	TTFBMedian time.Duration

	TotalTimePercentage []QuantileResult

	ErrorRequests,
	SuccessRequests,
	TotalRequests,
	Code2xx,
	Code3xx,
	Code4xx,
	Code5xx,
	OtherCodes int

	Server, PoweredBy string

	Errors []ErrorResult
}

type Statistics map[string]SingleUrlStatistics

func GetStatistics(results []tester.MeasurementResult, testDuration time.Duration) (Statistics, error) {
	var urls = make(map[string][]tester.MeasurementResult)

	for _, result := range results {
		url := result.RequestResult.Resource.Url.String()

		urls[url] = append(urls[url], result)
	}

	statistics := make(Statistics)

	for url, urlResults := range urls {
		stat, err := calculateStatistics(urlResults, testDuration)

		if err != nil {
			return nil, err
		}

		statistics[url] = stat
	}

	return statistics, nil
}

func calculateStatistics(results []tester.MeasurementResult, testDuration time.Duration) (SingleUrlStatistics, error) {
	var errorRequests, successRequests, totalRequests, code2xx, code3xx, code4xx, code5xx, otherCodes int
	var connectionEstablishedAvg, connectionEstablishedMin, connectionEstablishedMax time.Duration
	var tcpConnectionAvg, tcpConnectionMin, tcpConnectionMax time.Duration
	var tlsHandshakeAvg, tlsHandshakeMin, tlsHandshakeMax time.Duration
	var requestTimeAvg, requestTimeMin, requestTimeMax time.Duration
	var dnsLookupAvg, dnsLookupMin, dnsLookupMax time.Duration
	var ttfbAvg, ttfbMin, ttfbMax time.Duration

	errors := make(map[string]int)

	var timingPool struct {
		totalTime, dnsLookup, tcpConnection, tlsHandshake, connectionEstablished, ttfb []time.Duration
	}

	for _, result := range results {
		if result.Error == nil && result.RequestResult.Error == nil {
			requestTimeAvg += result.RequestResult.Durations.Total.Total
			requestTimeMin = minDuration(requestTimeMin, result.RequestResult.Durations.Total.Total)
			requestTimeMax = maxDuration(requestTimeMax, result.RequestResult.Durations.Total.Total)
			timingPool.totalTime = append(timingPool.totalTime, result.RequestResult.Durations.Total.Total)

			dnsLookupAvg += result.RequestResult.Durations.DNSLookup.Duration
			dnsLookupMin = minDuration(dnsLookupMin, result.RequestResult.Durations.DNSLookup.Duration)
			dnsLookupMax = maxDuration(dnsLookupMax, result.RequestResult.Durations.DNSLookup.Duration)
			timingPool.dnsLookup = append(timingPool.dnsLookup, result.RequestResult.Durations.DNSLookup.Duration)

			tcpConnectionAvg += result.RequestResult.Durations.TCPConnection.Duration
			tcpConnectionMax = maxDuration(tcpConnectionMax, result.RequestResult.Durations.TCPConnection.Duration)
			tcpConnectionMin = minDuration(tcpConnectionMin, result.RequestResult.Durations.TCPConnection.Duration)
			timingPool.tcpConnection = append(
				timingPool.tcpConnection, result.RequestResult.Durations.TCPConnection.Duration,
			)

			tlsHandshakeAvg += result.RequestResult.Durations.TLSHandshake.Duration
			tlsHandshakeMin = minDuration(tlsHandshakeMin, result.RequestResult.Durations.TLSHandshake.Duration)
			tlsHandshakeMax = maxDuration(tlsHandshakeMax, result.RequestResult.Durations.TLSHandshake.Duration)
			timingPool.tlsHandshake = append(
				timingPool.tlsHandshake, result.RequestResult.Durations.TLSHandshake.Duration,
			)

			connectionEstablishedAvg += result.RequestResult.Durations.ConnectionEstablishment.Duration
			connectionEstablishedMin = minDuration(
				connectionEstablishedMin, result.RequestResult.Durations.ConnectionEstablishment.Duration,
			)
			connectionEstablishedMax = maxDuration(
				connectionEstablishedMax, result.RequestResult.Durations.ConnectionEstablishment.Duration,
			)
			timingPool.connectionEstablished = append(
				timingPool.connectionEstablished, result.RequestResult.Durations.ConnectionEstablishment.Duration,
			)

			ttfbAvg += result.RequestResult.Durations.TTFB.Duration
			ttfbMin = minDuration(ttfbMin, result.RequestResult.Durations.TTFB.Duration)
			ttfbMax = maxDuration(ttfbMax, result.RequestResult.Durations.TTFB.Duration)
			timingPool.ttfb = append(timingPool.ttfb, result.RequestResult.Durations.TTFB.Duration)

			successRequests++

			if result.RequestResult.StatusCode >= 200 && result.RequestResult.StatusCode < 300 {
				code2xx++
			} else if result.RequestResult.StatusCode >= 300 && result.RequestResult.StatusCode < 400 {
				code3xx++
			} else if result.RequestResult.StatusCode >= 400 && result.RequestResult.StatusCode < 500 {
				code4xx++
			} else if result.RequestResult.StatusCode >= 500 && result.RequestResult.StatusCode < 600 {
				code5xx++
			} else {
				otherCodes++
			}
		} else {
			errorRequests++
		}

		if result.Error != nil {
			errors[result.Error.Error()]++
		}

		if result.RequestResult.Error != nil {
			errors[result.RequestResult.Error.Error()]++
		}

		totalRequests++
	}

	errorResult := make([]ErrorResult, 0, len(errors))

	for err, count := range errors {
		errorResult = append(
			errorResult, ErrorResult{
				Message: err,
				Count:   count,
			},
		)
	}

	server, poweredBy := "", ""

	if len(results) > 0 {
		server = results[0].RequestResult.Headers.Server
		poweredBy = results[0].RequestResult.Headers.PoweredBy
	}

	return SingleUrlStatistics{
		Server:              server,
		PoweredBy:           poweredBy,
		RequestTimeAvg:      requestTimeAvg / time.Duration(len(results)),
		RequestTimeMin:      requestTimeMin,
		RequestTimeMax:      requestTimeMax,
		RequestTimeMedian:   calculateDurationMedian(timingPool.totalTime),
		TotalTimePercentage: splitDataIntoSegments(timingPool.totalTime, 10),

		DNSLookupAvg:    dnsLookupAvg / time.Duration(len(results)),
		DNSLookupMin:    dnsLookupMin,
		DNSLookupMax:    dnsLookupMax,
		DNSLookupMedian: calculateDurationMedian(timingPool.dnsLookup),

		TCPConnectionAvg:    tcpConnectionAvg / time.Duration(len(results)),
		TCPConnectionMin:    tcpConnectionMin,
		TCPConnectionMax:    tcpConnectionMax,
		TCPConnectionMedian: calculateDurationMedian(timingPool.tcpConnection),

		TLSHandshakeAvg:    tlsHandshakeAvg / time.Duration(len(results)),
		TLSHandshakeMin:    tlsHandshakeMin,
		TLSHandshakeMax:    tlsHandshakeMax,
		TLSHandshakeMedian: calculateDurationMedian(timingPool.tlsHandshake),

		ConnectionEstablishedAvg:    connectionEstablishedAvg / time.Duration(len(results)),
		ConnectionEstablishedMin:    connectionEstablishedMin,
		ConnectionEstablishedMax:    connectionEstablishedMax,
		ConnectionEstablishedMedian: calculateDurationMedian(timingPool.connectionEstablished),

		TTFBAvg:    ttfbAvg / time.Duration(len(results)),
		TTFBMin:    ttfbMin,
		TTFBMax:    ttfbMax,
		TTFBMedian: calculateDurationMedian(timingPool.ttfb),

		TotalTime:       testDuration,
		ErrorRequests:   errorRequests,
		SuccessRequests: successRequests,
		TotalRequests:   totalRequests,
		Code2xx:         code2xx,
		Code3xx:         code3xx,
		Code4xx:         code4xx,
		Code5xx:         code5xx,
		OtherCodes:      otherCodes,

		Errors: errorResult,
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

func calculateDurationMedian(data []time.Duration) time.Duration {
	slices.Sort(data)

	n := len(data)

	if n == 0 {
		return 0
	}

	if n%2 == 1 {
		return data[n/2]
	}

	return (data[n/2-1] + data[n/2]) / 2
}

func splitDataIntoSegments(data []time.Duration, numSegments int) []QuantileResult {
	dataLength := len(data)

	if dataLength == 0 || numSegments <= 0 {
		return nil
	}

	slices.Sort(data)

	segmentSize := dataLength / numSegments
	segments := make([]QuantileResult, 0, numSegments)

	if dataLength < numSegments {
		segments = append(
			segments, QuantileResult{
				Segment: numSegments,
				Min:     data[0],
				Max:     data[dataLength-1],
			},
		)

		return segments
	}

	for i := 0; i < numSegments; i++ {
		startIndex := i * segmentSize
		endIndex := (i + 1) * segmentSize

		minValue := data[startIndex]
		maxValue := data[endIndex-1]

		segments = append(
			segments, QuantileResult{
				Segment: i + 1,
				Min:     minValue,
				Max:     maxValue,
			},
		)
	}

	return segments
}
