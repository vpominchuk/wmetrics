package tester

import (
	"context"
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strings"
	"time"
	"webmetrics/wmetrics/src/app"
)

func (engine *HttpEngine) Measure(parameters Parameters) (MeasurementsResult, error) {
	parameters.Method = strings.ToUpper(parameters.Method)

	return engine.request(parameters)
}

func (engine *HttpEngine) request(parameters Parameters) (MeasurementsResult, error) {
	request, err := engine.newRequest(parameters)
	client := engine.newClient(parameters, request)

	if err != nil {
		panic(err)
	}

	var result MeasurementsResult

	trace := engine.newClientTrace(&result)

	request = request.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	if userAgent, ok := app.DefaultUserAgents[parameters.UserAgentTemplate]; ok {
		request.Header.Set("user-agent", userAgent)
	} else {
		request.Header.Set("user-agent", parameters.UserAgent)
	}

	var response *http.Response
	response, err = client.Do(request)

	result.Timing.TotalTime = time.Now()

	defer response.Body.Close()

	if err != nil {
		return result, &ResponseError{
			Message: "Failed to read response",
			Err:     err,
		}
	}

	result.Status = response.Status
	result.StatusCode = response.StatusCode
	result.ContentLength = response.ContentLength

	if response.TLS != nil {
		result.TLS.UseTLS = true

		if response.TLS.Version == tls.VersionTLS12 {
			result.TLS.TLSVersion = "TLSv1.2"
		} else if response.TLS.Version == tls.VersionTLS13 {
			result.TLS.TLSVersion = "TLSv1.3"
		} else {
			result.TLS.TLSVersion = "UNKNOWN"
		}
	} else {
		result.TLS.UseTLS = false
	}

	if result.Timing.DNSStart.IsZero() {
		result.Timing.DNSStart = result.Timing.DNSEnd
	}

	if result.Timing.TLSHandshakeStart.IsZero() {
		result.Timing.TLSHandshakeStart = result.Timing.TCPConnect
		result.Timing.TLSHandshakeEnd = result.Timing.TCPConnect
	}

	engine.calculateDurations(&result)

	return result, nil
}

func (engine *HttpEngine) newClient(parameters Parameters, request *http.Request) *http.Client {
	proxyURL, _ := url.Parse(parameters.Proxy)

	var network string

	if parameters.IPv4Only {
		network = "tcp4"
	} else if parameters.IPv6Only {
		network = "tcp6"
	} else {
		network = "tcp4"
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, _, addr string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout:   parameters.Timeout * time.Millisecond,
				KeepAlive: parameters.IdleConnTimeout * time.Millisecond,
			}).DialContext(ctx, network, addr)
		},
		MaxIdleConns:          parameters.MaxIdleConnections,
		MaxIdleConnsPerHost:   0,
		MaxConnsPerHost:       0,
		IdleConnTimeout:       parameters.IdleConnTimeout,
		ResponseHeaderTimeout: parameters.Timeout,
		DisableKeepAlives:     !parameters.KeepAlive,
		TLSHandshakeTimeout:   parameters.TLSHandshakeTimeout,
		ExpectContinueTimeout: 1 * time.Second,
		ForceAttemptHTTP2:     true,
	}

	if parameters.Proxy != "" {
		transport.Proxy = http.ProxyURL(proxyURL)
	}

	if parameters.Url.Scheme == "https" {
		host, _, err := net.SplitHostPort(request.Host)

		if err != nil {
			host = request.Host
		}

		certificates, err := engine.readClientPemCertificate(parameters.ClientCertificateFile)

		if err != nil {
			panic(err)
		}

		transport.TLSClientConfig = &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: parameters.AllowInsecureSSL,
			Certificates:       certificates,
			MinVersion:         tls.VersionTLS12,
		}
	}

	return &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// do not follow redirects
			return http.ErrUseLastResponse
		},
	}
}

func (engine *HttpEngine) readClientPemCertificate(filename string) ([]tls.Certificate, error) {
	if filename == "" {
		return nil, nil
	}

	var (
		pkeyPem []byte
		certPem []byte
	)

	certFile, err := os.ReadFile(filename)

	if err != nil {
		return nil, &CertificateFileError{
			FileName: filename,
			Err:      err,
		}
	}

	for {
		block, rest := pem.Decode(certFile)

		if block == nil {
			break
		}

		certFile = rest

		if strings.HasSuffix(block.Type, "PRIVATE KEY") {
			pkeyPem = pem.EncodeToMemory(block)
		}

		if strings.HasSuffix(block.Type, "CERTIFICATE") {
			certPem = pem.EncodeToMemory(block)
		}
	}

	cert, err := tls.X509KeyPair(certPem, pkeyPem)

	if err != nil {
		log.Fatalf("unable to load client cert and key pair: %v", err)
	}

	return []tls.Certificate{cert}, nil
}

func (engine *HttpEngine) newRequest(parameters Parameters) (*http.Request, error) {
	var request *http.Request
	var err error

	switch parameters.Method {
	case http.MethodGet:
		request, err = http.NewRequest(http.MethodGet, parameters.Resource, nil)
	case http.MethodHead:
		request, err = http.NewRequest(http.MethodHead, parameters.Resource, nil)
	case http.MethodPost:
		// request, err = http.NewRequest(http.MethodPost, parameters.Resource, nil)
	case http.MethodPut:
		// request, err = http.NewRequest(http.MethodPut, parameters.Resource, nil)
	case http.MethodPatch:
		// request, err = http.NewRequest(http.MethodPatch, parameters.Resource, nil)
	default:
		request, err = http.NewRequest(http.MethodGet, parameters.Resource, nil)
	}

	return request, err
}

func (engine *HttpEngine) newClientTrace(result *MeasurementsResult) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			result.Timing.Start = time.Now()
		},
		GotConn:              func(_ httptrace.GotConnInfo) { result.Timing.ServerConnect = time.Now() },
		GotFirstResponseByte: func() { result.Timing.TTFB = time.Now() },
		DNSStart:             func(_ httptrace.DNSStartInfo) { result.Timing.DNSStart = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { result.Timing.DNSEnd = time.Now() },
		ConnectStart: func(_, _ string) {
			if result.Timing.DNSEnd.IsZero() {
				result.Timing.DNSEnd = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				panic(fmt.Sprintf("Unable to connect to host %v: %v", addr, err))
			}

			result.Timing.TCPConnect = time.Now()
		},
		TLSHandshakeStart: func() { result.Timing.TLSHandshakeStart = time.Now() },
		TLSHandshakeDone:  func(_ tls.ConnectionState, _ error) { result.Timing.TLSHandshakeEnd = time.Now() },
		WroteRequest:      func(_ httptrace.WroteRequestInfo) { result.Timing.RequestSent = time.Now() },
	}
}

func (engine *HttpEngine) calculateDurations(result *MeasurementsResult) {
	result.Durations.DNSLookup.Duration = result.Timing.DNSEnd.Sub(result.Timing.DNSStart)
	result.Durations.DNSLookup.Total = result.Timing.DNSEnd.Sub(result.Timing.Start)

	result.Durations.TCPConnection.Duration = result.Timing.TCPConnect.Sub(result.Timing.DNSEnd)
	result.Durations.TCPConnection.Total = result.Timing.TCPConnect.Sub(result.Timing.Start)

	result.Durations.TLSHandshake.Duration = result.Timing.TLSHandshakeEnd.Sub(result.Timing.TLSHandshakeStart)
	result.Durations.TLSHandshake.Total = result.Timing.TLSHandshakeEnd.Sub(result.Timing.Start)

	result.Durations.ConnectionEstablishment.Duration = result.Timing.ServerConnect.Sub(result.Timing.TLSHandshakeEnd)
	result.Durations.ConnectionEstablishment.Total = result.Timing.ServerConnect.Sub(result.Timing.Start)

	result.Durations.TTFB.Duration = result.Timing.TTFB.Sub(result.Timing.ServerConnect)
	result.Durations.TTFB.Total = result.Timing.TTFB.Sub(result.Timing.Start)

	result.Durations.Total.Duration = result.Timing.TotalTime.Sub(result.Timing.RequestSent)
	result.Durations.Total.Total = result.Timing.TotalTime.Sub(result.Timing.Start)
}
