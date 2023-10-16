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

type Arguments struct {
	Requests    intArgument
	Concurrency intArgument
	Timeout     durationArgument
	Method      stringArgument
	UserAgent   stringArgument
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
		name: "s", defaultValue: time.Millisecond * 30000,
		help: "`Milliseconds` to max. wait for each response.",
	},

	Method: stringArgument{
		name: "m", defaultValue: "GET",
		help: "HTTP `method`",
	},

	UserAgent: stringArgument{
		name: "u", defaultValue: app.UserAgent,
		help: "`User Agent`",
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
