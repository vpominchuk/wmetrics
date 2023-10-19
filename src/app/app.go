package app

var GitCommit string
var GitTag = "v0.0.1" // MAJOR.MINOR.PATCH
var ExecutableName string
var DefaultUserAgent = "wmetrics/" + GitTag
var VersionString = GitTag + " (build: " + GitCommit + ")"
var DefaultUserAgents = map[string]string{
	"chrome":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"chrome-mac":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"chrome-linux":   "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36",
	"chrome-android": "Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.5993.65 Mobile Safari/537.36",
	"chrome-iphone":  "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/118.0.5993.69 Mobile/15E148 Safari/604.1",
	"chrome-ipad":    "Mozilla/5.0 (iPad; CPU OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/118.0.5993.69 Mobile/15E148 Safari/604.1",
	"chrome-ipod":    "Mozilla/5.0 (iPod; CPU iPhone OS 17_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/118.0.5993.69 Mobile/15E148 Safari/604.1",

	"firefox":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/118.0",
	"firefox-mac":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 14.0; rv:109.0) Gecko/20100101 Firefox/118.0",
	"firefox-linux":   "Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/118.0",
	"firefox-android": "Mozilla/5.0 (Android 14; Mobile; rv:109.0) Gecko/118.0 Firefox/118.0",
	"firefox-iphone":  "Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/118.0 Mobile/15E148 Safari/605.1.15",
	"firefox-ipad":    "Mozilla/5.0 (iPad; CPU OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) FxiOS/118.0 Mobile/15E148 Safari/605.1.15",
	"firefox-ipod":    "Mozilla/5.0 (iPod touch; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/604.5.6 (KHTML, like Gecko) FxiOS/118.0 Mobile/15E148 Safari/605.1.15",

	"edge":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/117.0.2045.60",
	"edge-mac":     "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edg/117.0.2045.60",
	"edge-android": "Mozilla/5.0 (Linux; Android 10; HD1913) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.5993.65 Mobile Safari/537.36 EdgA/117.0.2045.65",
	"edge-iphone":  "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 EdgiOS/117.2045.65 Mobile/15E148 Safari/605.1.15",
	"edge-ipad":    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 EdgiOS/117.2045.65 Mobile/15E148 Safari/605.1.15",
	"edge-ipod":    "Mozilla/5.0 (iPhone; CPU iPhone OS 17_0_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 EdgiOS/117.2045.65 Mobile/15E148 Safari/605.1.15",
	"edge-xbox":    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; Xbox; Xbox One) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 Edge/44.18363.8131",
}
