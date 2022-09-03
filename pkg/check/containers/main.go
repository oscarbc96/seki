package containers

import "github.com/oscarbc96/seki/pkg/check"

var Checks []check.Check

func init() {
	Checks = []check.Check{
		new(AdviseDockerHubRateLimit),
		new(LatestTagIsNotUsed),
		new(PreferCopyOverAdd),
		new(LastUserIsNotRoot),
		new(WorkDirPathIsNotRelative),
	}
}
