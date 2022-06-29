package check

import "github.com/oscarbc96/seki/pkg/result"

type Run func() (*result.CheckResult, error)

var Analysers = make(map[string]Run)
