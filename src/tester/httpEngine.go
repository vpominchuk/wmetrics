package tester

import (
	"context"
	"crypto/tls"
	"encoding/pem"
	"github.com/vpominchuk/wmetrics/src/app"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type HttpEngine struct {
	Progress       RequestsProgress
	progressMutex  sync.Mutex
	resourceFeeder *ResourceFeeder
}

func (engine *HttpEngine) Measure(
	parameters Parameters,
	resourceFeeder *ResourceFeeder,
	onProgress func(progress RequestsProgress),
) ([]MeasurementResult, time.Duration) {
	engine.resourceFeeder = resourceFeeder

	engine.Progress = RequestsProgress{
		TotalRequests:     parameters.Requests,
		CompletedRequests: 0,
		FailedRequests:    0,
	}

	parameters.Method = strings.ToUpper(parameters.Method)

	results := make([]MeasurementResult, 0, parameters.Requests)
	concurrencyCh := make(chan bool, parameters.Concurrency)

	var wg sync.WaitGroup
	var lastOnProgressCalled int64 = 0
	var testStartTime time.Time

	testStartTime = time.Now()

	for requestNumber := 0; !engine.timeLimitReached(
		testStartTime, parameters.TimeLimit,
	) || requestNumber < parameters.Requests; requestNumber++ {
		if parameters.TimeLimit > 0 && engine.timeLimitReached(testStartTime, parameters.TimeLimit) {
			break
		}

		concurrencyCh <- true
		wg.Add(1)

		go func() {
			defer func() {
				<-concurrencyCh
				wg.Done()
			}()

			if parameters.TimeLimit > 0 && engine.timeLimitReached(testStartTime, parameters.TimeLimit) {
				return
			}

			result, err := engine.request(parameters)

			engine.updateProgress(err != nil)

			if time.Now().UnixNano()-lastOnProgressCalled >= int64(time.Second) {
				lastOnProgressCalled = time.Now().UnixNano()
				onProgress(engine.Progress)
			}

			results = append(
				results, MeasurementResult{
					RequestResult: result,
					Error:         err,
				},
			)
		}()
	}

	wg.Wait()
	close(concurrencyCh)

	onProgress(engine.Progress)
	return results, time.Now().Sub(testStartTime)
}

func (engine *HttpEngine) timeLimitReached(testStartTime time.Time, timeLimit time.Duration) bool {
	return time.Now().Sub(testStartTime) >= timeLimit
}

func (engine *HttpEngine) updateProgress(updateFailedRequests bool) {
	engine.progressMutex.Lock()

	engine.Progress.CompletedRequests++

	if updateFailedRequests {
		engine.Progress.FailedRequests++
	}

	engine.progressMutex.Unlock()
}

func (engine *HttpEngine) GetProgress() RequestsProgress {
	return engine.Progress
}

func (engine *HttpEngine) request(parameters Parameters) (RequestResult, error) {
	request, err := engine.newRequest(parameters)
	client := engine.newClient(parameters, request)

	if err != nil {
		return RequestResult{}, err
	}

	var result RequestResult

	trace := engine.newClientTrace(&result)
	request = request.WithContext(httptrace.WithClientTrace(context.Background(), trace))

	engine.setHeaders(parameters, request)

	var response *http.Response
	response, err = client.Do(request)

	result.Resource.Url = request.URL

	result.Timing.TotalTime = time.Now()

	engine.calculateDurations(&result)

	if err != nil {
		return result, &ResponseError{
			Message: "Failed to read response",
			Err:     err,
		}
	}

	defer response.Body.Close()

	result.Status = response.Status
	result.StatusCode = response.StatusCode
	result.ContentLength = response.ContentLength

	engine.fillTLSInfo(response, &result)
	engine.fillHeaders(response, &result)

	return result, nil
}

func (engine *HttpEngine) setHeaders(parameters Parameters, request *http.Request) {
	if userAgent, ok := app.DefaultUserAgents[parameters.UserAgentTemplate]; ok {
		request.Header.Set("user-agent", userAgent)
	} else {
		request.Header.Set("user-agent", parameters.UserAgent)
	}

	if parameters.ContentType != "" {
		request.Header.Set("content-type", parameters.ContentType)
	}

	if parameters.FormData != "" {
		request.Header.Set("content-type", "application/x-www-form-urlencoded")
	}

	if parameters.CustomHeaders != nil && len(parameters.CustomHeaders) > 0 {
		for _, header := range parameters.CustomHeaders {
			headerParts := strings.SplitN(header, ":", 2)

			request.Header.Set(headerParts[0], headerParts[1])
		}
	}
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

	if request.URL.Scheme == "https" {
		host, _, err := net.SplitHostPort(request.Host)

		if err != nil {
			host = request.Host
		}

		certificates, err := engine.readClientPemCertificate(parameters.ClientCertificateFile)

		if err != nil {
			log.Fatalf("Error: %v", err)
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

	resource, err := engine.resourceFeeder.GetNextValue()

	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	switch parameters.Method {
	case http.MethodGet, http.MethodHead, http.MethodDelete:
		request, err = http.NewRequest(parameters.Method, resource.Url.String(), nil)
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		postDataFileReader, err := engine.getPostDataReader(parameters)

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		request, err = http.NewRequest(parameters.Method, resource.Url.String(), postDataFileReader)
	default:
		request, err = http.NewRequest(http.MethodGet, resource.Url.String(), nil)
	}

	return request, err
}

func (engine *HttpEngine) getPostDataReader(parameters Parameters) (io.Reader, error) {
	if parameters.PostDataFile != "" {
		return engine.getPostDataFileReader(parameters.PostDataFile)
	}

	if parameters.FormData != "" {
		return strings.NewReader(parameters.FormData), nil
	}

	if parameters.PostData != "" {
		return strings.NewReader(parameters.PostData), nil
	}

	return nil, nil
}

func (engine *HttpEngine) getPostDataFileReader(filename string) (*os.File, error) {
	if filename == "" {
		return nil, nil
	}

	file, err := os.Open(filename)

	if err != nil {
		return nil, &PostDataFileError{
			FileName: filename,
			Err:      err,
		}
	}

	return file, nil
}

func (engine *HttpEngine) newClientTrace(result *RequestResult) *httptrace.ClientTrace {
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
			result.Error = err
			result.Timing.TCPConnect = time.Now()
		},
		TLSHandshakeStart: func() { result.Timing.TLSHandshakeStart = time.Now() },
		TLSHandshakeDone:  func(_ tls.ConnectionState, _ error) { result.Timing.TLSHandshakeEnd = time.Now() },
		WroteRequest:      func(_ httptrace.WroteRequestInfo) { result.Timing.RequestSent = time.Now() },
	}
}

func (engine *HttpEngine) calculateDurations(result *RequestResult) {
	if result.Timing.DNSStart.IsZero() {
		result.Timing.DNSStart = result.Timing.DNSEnd
	}

	if result.Timing.TLSHandshakeStart.IsZero() {
		result.Timing.TLSHandshakeStart = result.Timing.TCPConnect
		result.Timing.TLSHandshakeEnd = result.Timing.TCPConnect
	}

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

func (engine *HttpEngine) fillTLSInfo(response *http.Response, result *RequestResult) {
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
}

func (engine *HttpEngine) fillHeaders(response *http.Response, result *RequestResult) {
	result.Headers.Server = response.Header.Get("server")
	result.Headers.PoweredBy = response.Header.Get("x-powered-by")
}
