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

func TestLatestTagIsNotUsed(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    check.CheckResult
	}{
		{
			name: "Tag not specified",
			content: []byte(`FROM ubuntu
COPY . /app
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  LatestTagIsNotUsed{},
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
			name: "Tag is latest",
			content: []byte(`FROM ubuntu:latest
COPY . /app
USER root
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  LatestTagIsNotUsed{},
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
			name: "Tag is not latest",
			content: []byte(`FROM ubuntu:18.04
COPY . /app
USER potatoe
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  LatestTagIsNotUsed{},
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

			c := new(LatestTagIsNotUsed)
			assert.Equal(t, tt.want, c.Run(*load.NewInput("Dockerfile", nil)))
		})
	}
}
