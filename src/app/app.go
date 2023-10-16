package app

var GitCommit string
var GitTag = "v0.0.1" // MAJOR.MINOR.PATCH
var ExecutableName string
var UserAgent = "wmetrics/" + GitTag
var VersionString = GitTag + " (build: " + GitCommit + ")"
