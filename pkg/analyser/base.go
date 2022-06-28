package analyser

import "github.com/oscarbc96/seki/pkg/result"

type Run func() (*result.RuleResult, error)

var Analysers = make(map[string]Run)
