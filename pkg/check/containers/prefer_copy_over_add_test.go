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

func TestPreferCopyOverAdd(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    check.CheckResult
	}{
		{
			name: "ADD not specified",
			content: []byte(`FROM ubuntu:18.04
COPY . /app
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  PreferCopyOverAdd{},
				Status: check.PASS,
			},
		},
		{
			name: "ADD is present",
			content: []byte(`FROM ubuntu:latest
ADD . /app
USER root
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  PreferCopyOverAdd{},
				Status: check.FAIL,
				Locations: check.Locations{
					load.Range{
						Start: load.Position{Line: 2, Column: 1},
						End:   load.Position{Line: 2, Column: 3},
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

			c := new(PreferCopyOverAdd)
			assert.Equal(t, tt.want, c.Run(*load.NewInput("Dockerfile", nil)))
		})
	}
}
