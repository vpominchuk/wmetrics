package tester

import (
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
}

type TestEngine interface {
	Measure(parameters Parameters) (MeasurementsResult, error)
}

type HttpEngine struct {
}

type Timing struct {
	dnsStart,
	dnsEnd,
	tcpConnect,
	serverConnect,
	ttfb,
	tlsHandshakeStart,
	tlsHandshakeEnd,
	requestSent,
	totalTime time.Time
}

type TLS struct {
	UseTLS     bool
	TLSVersion string
}

type MeasurementsResult struct {
	Status        string // e.g. "200 OK"
	StatusCode    int    // e.g. 200
	ContentLength int64
	Timing        Timing
	TLS           TLS
}
