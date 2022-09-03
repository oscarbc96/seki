package cloudformation

import "github.com/oscarbc96/seki/pkg/check"

var AllChecks []check.Check

func init() {
	AllChecks = append(
		AllChecks,
		new(S3BucketHasAllBlockPublicAccessEnabled),
		new(S3BucketVersioningIsEnabled),
	)
}
