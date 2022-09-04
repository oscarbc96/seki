package metadata

import (
	"runtime"
	"time"
)

var Version = "development"
var Commit = "development"
var BuiltDate = time.Now().UTC().Format(time.RFC3339)
var RuntimeOS = runtime.GOOS
var RuntimeArch = runtime.GOARCH
var Homepage = "https://oscarbc96.github.io/seki/"
