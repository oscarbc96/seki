package metadata

import (
	"fmt"
	"runtime"
	"time"
)

var Version = "development"
var Commit = "development"
var BuiltDate = time.Now().UTC().Format(time.RFC3339)
var RuntimeOS = runtime.GOOS
var RuntimeArch = runtime.GOARCH
var Homepage = "https://oscarbc96.github.io/seki/"

func GenerateChecksDocsURL(path string) string {
	return fmt.Sprintf("%sdocs/checks/%s", Homepage, path)
}
