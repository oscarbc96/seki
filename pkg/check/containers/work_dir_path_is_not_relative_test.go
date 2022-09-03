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

func TestWorkDirPathIsNotRelative(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    check.CheckResult
	}{
		{
			name: "WORKDIR is absolute",
			content: []byte(`FROM ubuntu:18.04
WORKDIR /app
COPY . /app
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  WorkDirPathIsNotRelative{},
				Status: check.PASS,
			},
		},
		{
			name: "WORKDIR is relative",
			content: []byte(`FROM ubuntu:latest
WORKDIR app
ADD . /app
USER root
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  WorkDirPathIsNotRelative{},
				Status: check.FAIL,
				Locations: check.Locations{
					load.Range{
						Start: load.Position{Line: 2, Column: 9},
						End:   load.Position{Line: 2, Column: 12},
					},
				},
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

			c := new(WorkDirPathIsNotRelative)
			assert.Equal(t, tt.want, c.Run(*load.NewInput("Dockerfile", nil)))
		})
	}
}
