package args

import (
	"flag"
	"fmt"
	"time"
	"webmetrics/wmetrics/src/app"
)

type stringArgument struct {
	name         string
	help         string
	defaultValue string
	Value        *string
}

type intArgument struct {
	name         string
	help         string
	defaultValue int
	Value        *int
}

type durationArgument struct {
	name         string
	help         string
	defaultValue time.Duration
	Value        *time.Duration
}

type boolArgument struct {
	name         string
	help         string
	defaultValue bool
	Value        *bool
}

type Arguments struct {
	Requests              intArgument
	Concurrency           intArgument
	Timeout               durationArgument
	Method                stringArgument
	UserAgent             stringArgument
	UserAgentTemplate     stringArgument
	KeepAlive             boolArgument
	Proxy                 stringArgument
	MaxIdleConnections    intArgument
	IdleConnTimeout       durationArgument
	TLSHandshakeTimeout   durationArgument
	IPv4Only              boolArgument
	IPv6Only              boolArgument
	AllowInsecureSSL      boolArgument
	ClientCertificateFile stringArgument
}

var arguments = Arguments{
	Requests: intArgument{
		name: "n", defaultValue: 1,
		help: "Number of `requests` to perform",
	},

	Concurrency: intArgument{
		name: "c", defaultValue: 1,
		help: "Number of multiple `requests` to make at a time",
	},

	Timeout: durationArgument{
		name: "s", defaultValue: 30 * time.Second,
		help: "`Milliseconds` to max. wait for each response.",
	},

	Method: stringArgument{
		name: "m", defaultValue: "GET",
		help: "HTTP `method`",
	},

	UserAgent: stringArgument{
		name: "u", defaultValue: app.DefaultUserAgent,
		help: "`User Agent`",
	},

	UserAgentTemplate: stringArgument{
		name: "ut", defaultValue: "",
		help: "Use `User Agent Template`. Allowed values (chrome, firefox, edge)[-(linux, mac, android, iphone, ipod, ipad)]",
	},

	KeepAlive: boolArgument{
		name: "k", defaultValue: false,
		help: "Use HTTP KeepAlive feature",
	},

	MaxIdleConnections: intArgument{
		name: "km", defaultValue: 100,
		help: "Max idle `connections`",
	},

	IdleConnTimeout: durationArgument{
		name: "kt", defaultValue: 90 * time.Second,
		help: "Max idle connections `timeout` in ms",
	},

	Proxy: stringArgument{
		name: "P", defaultValue: "",
		help: "Use proxy. Values may be either a complete `URL` or a \"host[:port]\". The schemes \"http\", \"https\", and \"socks5\" are supported.",
	},

	TLSHandshakeTimeout: durationArgument{
		name: "tt", defaultValue: 10 * time.Second,
		help: "TLS handshake `timeout` in ms",
	},

	IPv4Only: boolArgument{
		name: "4", defaultValue: true,
		help: "Resolve IPv4 addresses only",
	},

	IPv6Only: boolArgument{
		name: "6", defaultValue: false,
		help: "Resolve IPv6 addresses only",
	},

	AllowInsecureSSL: boolArgument{
		name: "i", defaultValue: false,
		help: "Allow insecure SSL connections",
	},

	ClientCertificateFile: stringArgument{
		name: "C", defaultValue: "",
		help: "Client PEM certificate `file`",
	},
}

func (arguments *Arguments) init() {
	arguments.Requests.Value = flag.Int(
		arguments.Requests.name, arguments.Requests.defaultValue, arguments.Requests.help,
	)

	arguments.Concurrency.Value = flag.Int(
		arguments.Concurrency.name, arguments.Concurrency.defaultValue, arguments.Concurrency.help,
	)

	arguments.Timeout.Value = flag.Duration(
		arguments.Timeout.name, arguments.Timeout.defaultValue, arguments.Timeout.help,
	)

	arguments.Method.Value = flag.String(
		arguments.Method.name, arguments.Method.defaultValue, arguments.Method.help,
	)

	arguments.UserAgent.Value = flag.String(
		arguments.UserAgent.name, arguments.UserAgent.defaultValue, arguments.UserAgent.help,
	)

	arguments.UserAgentTemplate.Value = flag.String(
		arguments.UserAgentTemplate.name, arguments.UserAgentTemplate.defaultValue, arguments.UserAgentTemplate.help,
	)

	arguments.KeepAlive.Value = flag.Bool(
		arguments.KeepAlive.name, arguments.KeepAlive.defaultValue, arguments.KeepAlive.help,
	)

	arguments.Proxy.Value = flag.String(
		arguments.Proxy.name, arguments.Proxy.defaultValue, arguments.Proxy.help,
	)

	arguments.MaxIdleConnections.Value = flag.Int(
		arguments.MaxIdleConnections.name, arguments.MaxIdleConnections.defaultValue, arguments.MaxIdleConnections.help,
	)

	arguments.IdleConnTimeout.Value = flag.Duration(
		arguments.IdleConnTimeout.name, arguments.IdleConnTimeout.defaultValue, arguments.IdleConnTimeout.help,
	)

	arguments.TLSHandshakeTimeout.Value = flag.Duration(
		arguments.TLSHandshakeTimeout.name, arguments.TLSHandshakeTimeout.defaultValue,
		arguments.TLSHandshakeTimeout.help,
	)

	arguments.IPv4Only.Value = flag.Bool(
		arguments.IPv4Only.name, arguments.IPv4Only.defaultValue, arguments.IPv4Only.help,
	)

	arguments.IPv6Only.Value = flag.Bool(
		arguments.IPv6Only.name, arguments.IPv6Only.defaultValue, arguments.IPv6Only.help,
	)

	arguments.AllowInsecureSSL.Value = flag.Bool(
		arguments.AllowInsecureSSL.name, arguments.AllowInsecureSSL.defaultValue, arguments.AllowInsecureSSL.help,
	)

	arguments.ClientCertificateFile.Value = flag.String(
		arguments.ClientCertificateFile.name, arguments.ClientCertificateFile.defaultValue,
		arguments.ClientCertificateFile.help,
	)

	flag.Parse()
}

func GetArguments() (Arguments, []string) {
	arguments.init()
	return arguments, flag.Args()
}

func Usage() {
	fmt.Fprintf(
		flag.CommandLine.Output(), "Usage: %s [options] http[s]://]hostname[:port][/path]\n", app.ExecutableName,
	)
	fmt.Fprint(flag.CommandLine.Output(), "Options are:\n")
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nVersion: %s\n", app.VersionString)
}

func Validate(arguments Arguments) {

}
