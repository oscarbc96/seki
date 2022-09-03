package containers

import (
	"github.com/oscarbc96/seki/pkg/check"
	"github.com/oscarbc96/seki/pkg/load"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"testing"
)

func init() {
	load.AFS = &afero.Afero{Fs: afero.NewMemMapFs()}
}

func TestAdviseDockerHubRateLimitRun(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    check.CheckResult
	}{
		{
			name: "Registry not specified",
			content: []byte(`FROM ubuntu:18.04
COPY . /app
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  AdviseDockerHubRateLimit{},
				Status: check.FAIL,
				Locations: check.Locations{
					load.Range{
						Start: load.Position{Line: 1, Column: 0},
						End:   load.Position{Line: 1, Column: 0},
					},
				},
			},
		},
		{
			name: "Registry is docker hub",
			content: []byte(`FROM docker.io/ubuntu:18.04
COPY . /app
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  AdviseDockerHubRateLimit{},
				Status: check.q,
				Locations: check.Locations{
					load.Range{
						Start: load.Position{Line: 1, Column: 0},
						End:   load.Position{Line: 1, Column: 0},
					},
				},
			},
		},
		{
			name: "Registry is not docker hub",
			content: []byte(`FROM ghcr.io/ubuntu:18.04
COPY . /app
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  AdviseDockerHubRateLimit{},
				Status: check.PASS,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := load.AFS.WriteFile(
				"Dockerfile",
				tt.content,
				0644,
			)
			if err != nil {
				log.Fatal().Err(err)
			}

			c := new(AdviseDockerHubRateLimit)
			assert.Equal(t, tt.want, c.Run(*load.NewInput("Dockerfile", nil)))
		})
	}
}
