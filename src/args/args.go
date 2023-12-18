package args

import (
	"flag"
	"fmt"
	"github.com/vpominchuk/wmetrics/src/app"
	"github.com/vpominchuk/wmetrics/src/formatter"
	"strings"
	"time"
)

type stringArrayArgument struct {
	Name         string
	help         string
	defaultValue []string
	Value        *[]string
}

type stringArgument struct {
	Name         string
	help         string
	defaultValue string
	Value        *string
}

type intArgument struct {
	Name         string
	help         string
	defaultValue int
	Value        *int
}

type durationArgument struct {
	Name         string
	help         string
	defaultValue time.Duration
	Value        *time.Duration
}

type boolArgument struct {
	Name         string
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
	PostDataFile          stringArgument
	PostData              stringArgument
	ContentType           stringArgument
	FormData              stringArgument
	OutputFormat          stringArgument
	CustomHeaders         stringArrayArgument
	TimeLimit             durationArgument
	URLListFile           stringArgument
	ExitWithErrorOnCode   stringArrayArgument
}

type multipleStringValues []string

func (s *multipleStringValues) String() string {
	return strings.Join(*s, "\n")
}

func (s *multipleStringValues) Set(value string) error {
	*s = append(*s, value)
	return nil
}

var arguments = Arguments{
	Requests: intArgument{
		Name: "n", defaultValue: 1,
		help: "Number of `requests` to perform",
	},

	Concurrency: intArgument{
		Name: "c", defaultValue: 1,
		help: "Number of multiple `requests` to make at a time",
	},

	Timeout: durationArgument{
		Name: "s", defaultValue: 30 * time.Second,
		help: "`time` (30s, 800ms, ...) to max. wait for each response",
	},

	Method: stringArgument{
		Name: "m", defaultValue: "GET",
		help: "HTTP `method`",
	},

	UserAgent: stringArgument{
		Name: "u", defaultValue: app.DefaultUserAgent,
		help: "`User Agent`",
	},

	UserAgentTemplate: stringArgument{
		Name: "ut", defaultValue: "",
		help: "Use `User Agent Template`. Allowed values (chrome, firefox, edge)[-(linux, mac, android, iphone, ipod, ipad)] [-ut list] to see all templates",
	},

	KeepAlive: boolArgument{
		Name: "k", defaultValue: false,
		help: "Use HTTP KeepAlive feature",
	},

	MaxIdleConnections: intArgument{
		Name: "km", defaultValue: 100,
		help: "Max idle `connections`",
	},

	IdleConnTimeout: durationArgument{
		Name: "kt", defaultValue: 90 * time.Second,
		help: "Max idle connections `timeout` (90s, 800ms, ...)",
	},

	Proxy: stringArgument{
		Name: "P", defaultValue: "",
		help: "Use proxy. Values may be either a complete `URL` or a \"host[:port]\". The schemes \"http\", \"https\", and \"socks5\" are supported.",
	},

	TLSHandshakeTimeout: durationArgument{
		Name: "tt", defaultValue: 10 * time.Second,
		help: "TLS handshake `timeout` (10s, 800ms, ...)",
	},

	IPv4Only: boolArgument{
		Name: "4", defaultValue: true,
		help: "Resolve IPv4 addresses only",
	},

	IPv6Only: boolArgument{
		Name: "6", defaultValue: false,
		help: "Resolve IPv6 addresses only",
	},

	AllowInsecureSSL: boolArgument{
		Name: "i", defaultValue: false,
		help: "Allow insecure SSL connections",
	},

	ClientCertificateFile: stringArgument{
		Name: "C", defaultValue: "",
		help: "Client PEM certificate `file`",
	},

	PostDataFile: stringArgument{
		Name: "f", defaultValue: "",
		help: "Post data `file`",
	},

	PostData: stringArgument{
		Name: "d", defaultValue: "",
		help: "Post data `string`",
	},

	ContentType: stringArgument{
		Name: "T", defaultValue: "application/json",
		help: "Content type",
	},

	FormData: stringArgument{
		Name: "F", defaultValue: "",
		help: "Form data",
	},

	OutputFormat: stringArgument{
		Name: "O", defaultValue: "std",
		help: "Output `format`. Allowed values (std, text, json, json-pretty)",
	},

	CustomHeaders: stringArrayArgument{
		Name: "H", defaultValue: nil,
		help: "Custom `header`. For example: \"Accept-Encoding: gzip, deflate\"." +
			" Multiple headers can be provided with multiple -H flags.",
	},

	TimeLimit: durationArgument{
		Name: "t", defaultValue: 0,
		help: "`Time` limit (1s, 200ms, ...). If the time limit is reached, " + app.ExecutableName +
			" will interrupt the test and print the results.",
	},

	URLListFile: stringArgument{
		Name: "l", defaultValue: "",
		help: "Path to a `file` with list of URLs",
	},

	ExitWithErrorOnCode: stringArrayArgument{
		Name: "e", defaultValue: nil,
		help: "Exit with error on HTTP `code`. Multiple codes can be provided with multiple -e flags. Example: -e 403 -e 3xx.",
	},
}

func (arguments *Arguments) init() {
	arguments.Requests.Value = flag.Int(
		arguments.Requests.Name, arguments.Requests.defaultValue, arguments.Requests.help,
	)

	arguments.Concurrency.Value = flag.Int(
		arguments.Concurrency.Name, arguments.Concurrency.defaultValue, arguments.Concurrency.help,
	)

	arguments.Timeout.Value = flag.Duration(
		arguments.Timeout.Name, arguments.Timeout.defaultValue, arguments.Timeout.help,
	)

	arguments.Method.Value = flag.String(
		arguments.Method.Name, arguments.Method.defaultValue, arguments.Method.help,
	)

	arguments.UserAgent.Value = flag.String(
		arguments.UserAgent.Name, arguments.UserAgent.defaultValue, arguments.UserAgent.help,
	)

	arguments.UserAgentTemplate.Value = flag.String(
		arguments.UserAgentTemplate.Name, arguments.UserAgentTemplate.defaultValue, arguments.UserAgentTemplate.help,
	)

	arguments.KeepAlive.Value = flag.Bool(
		arguments.KeepAlive.Name, arguments.KeepAlive.defaultValue, arguments.KeepAlive.help,
	)

	arguments.Proxy.Value = flag.String(
		arguments.Proxy.Name, arguments.Proxy.defaultValue, arguments.Proxy.help,
	)

	arguments.MaxIdleConnections.Value = flag.Int(
		arguments.MaxIdleConnections.Name, arguments.MaxIdleConnections.defaultValue, arguments.MaxIdleConnections.help,
	)

	arguments.IdleConnTimeout.Value = flag.Duration(
		arguments.IdleConnTimeout.Name, arguments.IdleConnTimeout.defaultValue, arguments.IdleConnTimeout.help,
	)

	arguments.TLSHandshakeTimeout.Value = flag.Duration(
		arguments.TLSHandshakeTimeout.Name, arguments.TLSHandshakeTimeout.defaultValue,
		arguments.TLSHandshakeTimeout.help,
	)

	arguments.IPv4Only.Value = flag.Bool(
		arguments.IPv4Only.Name, arguments.IPv4Only.defaultValue, arguments.IPv4Only.help,
	)

	arguments.IPv6Only.Value = flag.Bool(
		arguments.IPv6Only.Name, arguments.IPv6Only.defaultValue, arguments.IPv6Only.help,
	)

	arguments.AllowInsecureSSL.Value = flag.Bool(
		arguments.AllowInsecureSSL.Name, arguments.AllowInsecureSSL.defaultValue, arguments.AllowInsecureSSL.help,
	)

	arguments.ClientCertificateFile.Value = flag.String(
		arguments.ClientCertificateFile.Name, arguments.ClientCertificateFile.defaultValue,
		arguments.ClientCertificateFile.help,
	)

	arguments.PostDataFile.Value = flag.String(
		arguments.PostDataFile.Name, arguments.PostDataFile.defaultValue,
		arguments.PostDataFile.help,
	)

	arguments.PostData.Value = flag.String(
		arguments.PostData.Name, arguments.PostData.defaultValue,
		arguments.PostData.help,
	)

	arguments.ContentType.Value = flag.String(
		arguments.ContentType.Name, arguments.ContentType.defaultValue,
		arguments.ContentType.help,
	)

	arguments.FormData.Value = flag.String(
		arguments.FormData.Name, arguments.FormData.defaultValue,
		arguments.FormData.help,
	)

	arguments.OutputFormat.Value = flag.String(
		arguments.OutputFormat.Name, arguments.OutputFormat.defaultValue,
		arguments.OutputFormat.help,
	)

	arguments.URLListFile.Value = flag.String(
		arguments.URLListFile.Name, arguments.URLListFile.defaultValue,
		arguments.URLListFile.help,
	)

	var headers multipleStringValues
	flag.Var(&headers, arguments.CustomHeaders.Name, arguments.CustomHeaders.help)
	arguments.CustomHeaders.Value = (*[]string)(&headers)

	arguments.TimeLimit.Value = flag.Duration(
		arguments.TimeLimit.Name, arguments.TimeLimit.defaultValue,
		arguments.TimeLimit.help,
	)

	var exitWithErrorOnCode multipleStringValues
	flag.Var(&exitWithErrorOnCode, arguments.ExitWithErrorOnCode.Name, arguments.ExitWithErrorOnCode.help)
	arguments.ExitWithErrorOnCode.Value = (*[]string)(&exitWithErrorOnCode)

	flag.Usage = customUsage

	flag.Parse()
}

func GetArguments() (Arguments, []string) {
	arguments.init()
	return arguments, flag.Args()
}

func Usage() {
	flag.Usage()
}

func customUsage() {
	fmt.Fprintf(
		flag.CommandLine.Output(), "Usage: %s [options] URL_LIST\n", app.ExecutableName,
	)
	fmt.Fprint(flag.CommandLine.Output(), "Options are:\n")
	customPrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), "\nVersion: %s\n", app.VersionString)
}

func customPrintDefaults() {
	strLength := 30
	flag.VisitAll(
		func(f *flag.Flag) {
			name, usage := flag.UnquoteUsage(f)

			usageMessage := formatter.StrPadRight(fmt.Sprintf("  -%s %s", f.Name, name), strLength)
			usageMessage += fmt.Sprintf(" %s", usage)

			if f.DefValue != "" {
				usageMessage += fmt.Sprintf(" (default: %v)", f.DefValue)
			}

			fmt.Println(usageMessage)
		},
	)
}
