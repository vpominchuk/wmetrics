package tester

import (
	"encoding/json"
	"net/url"
	"time"
)

type Parameters struct {
	Resource              string
	Url                   *url.URL
	Requests              int
	Concurrency           int
	Timeout               time.Duration
	Method                string
	UserAgent             string
	UserAgentTemplate     string
	KeepAlive             bool
	Proxy                 string
	MaxIdleConnections    int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	IPv4Only              bool
	IPv6Only              bool
	AllowInsecureSSL      bool
	ClientCertificateFile string
	PostDataFile          string
	PostData              string
	ContentType           string
	FormData              string
	OutputFormat          string
}

type TestEngine interface {
	Measure(parameters Parameters, onProgress func(progress RequestsProgress)) ([]MeasurementResult, time.Duration)
	GetProgress() RequestsProgress
}

type RequestsProgress struct {
	TotalRequests,
	CompletedRequests,
	FailedRequests int
}

type Timing struct {
	Start,
	DNSStart,
	DNSEnd,
	TCPConnect,
	ServerConnect,
	TTFB,
	TLSHandshakeStart,
	TLSHandshakeEnd,
	RequestSent,
	TotalTime time.Time
}

type Duration struct {
	Duration,
	Total time.Duration
}

type Durations struct {
	DNSLookup,
	TCPConnection,
	TLSHandshake,
	ConnectionEstablishment,
	TTFB,
	Total Duration
}

type TLS struct {
	UseTLS     bool
	TLSVersion string
}

type ResponseHeaders struct {
	Server,
	PoweredBy string
}

type RequestResult struct {
	Status        string // e.g. "200 OK"
	StatusCode    int    // e.g. 200
	ContentLength int64
	Timing        Timing
	Durations     Durations
	TLS           TLS
	Headers       ResponseHeaders
	Error         error
}

func (result RequestResult) ToJson() ([]byte, error) {
	return json.MarshalIndent(result, "", "  ")
}

type MeasurementResult struct {
	RequestResult RequestResult
	Error         error
}
