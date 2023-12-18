# Web Metrics

`wmetrics` is a versatile command-line tool written in Go for benchmarking and monitoring web servers. It allows you to simulate various HTTP requests, measure server performance, and collect valuable metrics. Whether you need to stress test your web application, analyze response times, or monitor your server's behavior, wmetrics is a handy tool to have in your toolkit.

## Features

- Send HTTP requests to a specified URL.
- Support for GET and POST requests with customizable data and content type.
- Control the number of concurrent requests to emulate real-world scenarios.
- Resolve IPv4 or IPv6 addresses as per your needs.
- Utilize client PEM certificates for secure connections.
- Support for proxy servers (HTTP, HTTPS, and SOCKS5).
- Fine-tune SSL and TLS settings.
- Specify custom User Agent or choose from predefined templates.
- Output data in multiple formats, including standard output, plain text, and JSON.

## Installation

To install `wmetrics`, you need to have Go installed on your system. Then, run the following command:

```bash
git clone https://github.com/vpominchuk/wmetrics
cd wmetrics
go install
./build
```

## Pre-built Binaries
You can download pre-built binaries for Linux, macOS, and Windows from the [Releases](https://github.com/vpominchuk/wmetrics/releases) page.


## Build
To build the project for different platforms and architectures please follow our [Build Guidelines](docs/BUILD.md).

## Usage
```bash
wmetrics [options] URL_LIST
```

## Options
| Option                  | Description                                                                                                                                     |
|-------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------|
| -4                      | Resolve IPv4 addresses only (default true).                                                                                                     |
| -6                      | Resolve IPv6 addresses only.                                                                                                                    |
| -C file                 | Client PEM certificate file.                                                                                                                    |
| -F string               | Form data.                                                                                                                                      |
| -H header               | Custom header. For example: "Accept-Encoding: gzip, deflate". Multiple headers can be provided with multiple -H flags.                          |
| -O format               | Output format (std, text, json) (default "std").                                                                                                |
| -P URL                  | Use a proxy (complete URL or "host[:port]"). Supported schemes: "http," "https," and "socks5."                                                  |
| -T string               | Content type (default "application/json").                                                                                                      |
| -c requests             | Number of concurrent requests (default 1).                                                                                                      |
| -d string               | Post data as a string.                                                                                                                          |
| -e code                 | Exit with error on HTTP code. Multiple code can be provided with multiple -e flags.                                                             |
| -f file                 | Post data from a file.                                                                                                                          |
| -i                      | Allow insecure SSL connections.                                                                                                                 |
| -k                      | Use HTTP KeepAlive feature.                                                                                                                     |
| -l file                 | URL list file                                                                                                                                   |
| -km connections         | Max idle connections (default 100).                                                                                                             |
| -kt timeout             | Max idle connections timeout in ms (default 1m30s).                                                                                             |
| -m method               | HTTP method (default "GET").                                                                                                                    |
| -n requests             | Number of requests to perform (default 1).                                                                                                      |
| -s Milliseconds         | Maximum wait time for each response (default 30s).                                                                                              |
| -t time                 | Time limit (1s, 200ms, ...). If the time limit is reached, wmetrics will interrupt the test and print the results.                              |
| -tt timeout             | TLS handshake timeout in ms (default 10s).                                                                                                      |
| -u User Agent           | User Agent (default "wmetrics/v0.0.1").                                                                                                         |
| -ut User Agent Template | Use User Agent Template. Allowed values (chrome, firefox, edge)[-(linux, mac, android, iphone, ipod, ipad)]. Use -ut list to see all templates. |


## Examples

### Send a GET request to a URL
```bash
wmetrics https://example.com
```

### Send a POST request to a URL
```bash
wmetrics -m POST -d '{"name":"John Doe"}' https://example.com
```

### Send a POST request with form data
```bash
wmetrics -m POST -F "first_name=John&last_name=Doe" https://example.com
```

### Send a POST request with data from a file
```bash
wmetrics -m POST -f data.json https://example.com
```

### Send a POST request with a client certificate
```bash
wmetrics -m POST -C client.pem -d '{"name":"John Doe"}' https://example.com
```

### Send requests with custom User Agent and output as JSON
```bash
wmetrics -u "MyUserAgent" -O json https://example.com
```

### Send requests with a predefined User Agent template
```bash
wmetrics -ut chrome-linux https://example.com
```

### Stress test a server with multiple concurrent connections
```bash
wmetrics -c 100 -n 1000 https://example.com
```

For more options and detailed usage, please refer to the program's help documentation.

## License
This project is licensed under the [MIT License](MIT-LICENSE.txt).

## Contributing
Contributions are welcome! If you would like to contribute to this project, please follow our [Contributing Guidelines](docs/CONTRIBUTING.md).