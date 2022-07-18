package check

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckIdsAreUnique(t *testing.T) {
	ids := lo.Map[Check, string](allChecks, func(c Check, _ int) string {
		return c.Id()
	})
	duplicatedIds := lo.FindDuplicates[string](ids)

	assert.Empty(t, duplicatedIds, "Found %v duplicated", duplicatedIds)
}
