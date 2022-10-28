package load

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	AFS = &afero.Afero{Fs: afero.NewMemMapFs()}
}

func TestDetectorCICDGitHubActions_GitHubWorkflowsPathRegex(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Extension .yml",
			path:     ".github/workflows/test.yml",
			expected: true,
		},
		{
			name:     "Extension .yaml",
			path:     ".github/workflows/test.yaml",
			expected: true,
		},
		{
			name:     "Workflow name lowercase",
			path:     ".github/workflows/test.yaml",
			expected: true,
		},
		{
			name:     "Workflow name uppercase",
			path:     ".github/workflows/TEST.yaml",
			expected: true,
		},
		{
			name:     "Workflow name with numbers",
			path:     ".github/workflows/0123456789.yaml",
			expected: true,
		},
		{
			name:     "Workflow name with hyphens",
			path:     ".github/workflows/aa-bb.yaml",
			expected: true,
		},
		{
			name:     "Workflow name with underscores",
			path:     ".github/workflows/aa_bb.yaml",
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GitHubWorkflowsPathRegex.MatchString(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectorCICDGitHubActions_Detect(t *testing.T) {
	tests := []struct {
		name        string
		content     []byte
		expected    DetectedType
		expectedErr error
	}{
		{
			name: "Basic workflow",
			content: []byte(`name: 'Dependency Review'
on: pull_request

permissions:
  contents: read

jobs:
  dependency-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0

      - name: 'Dependency Review'
        uses: actions/dependency-review-action@v2
`),
			expected:    DetectedCICDGitHubActions,
			expectedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AFS.WriteFile(
				"/.github/workflows/test.yml",
				tt.content,
				0644,
			)
			if err != nil {
				log.Fatal().Err(err)
			}

			de := DetectorCICDGitHubActions{}
			got, err := de.Detect(*NewInput("/.github/actions/test.yml", nil))
			assert.Equal(t, tt.expected, got)
			assert.ErrorIs(t, tt.expectedErr, err)
		})
	}
}
