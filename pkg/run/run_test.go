package run

import (
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckIdsAreUnique(t *testing.T) {
	ids := lo.Map[check.Check, string](AllChecks, func(c check.Check, _ int) string {
		return c.Id()
	})
	duplicatedIds := lo.FindDuplicates[string](ids)

	assert.Empty(t, duplicatedIds, "Found %v duplicated", duplicatedIds)
}

func TestNumberOfChecks(t *testing.T) {
	assert.Equal(t, 6, len(AllChecks))
}
