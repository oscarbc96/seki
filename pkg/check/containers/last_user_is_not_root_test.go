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

func TestLastUserIsNotRoot(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    check.CheckResult
	}{
		{
			name: "User not specified",
			content: []byte(`FROM ubuntu:18.04
COPY . /app
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  LastUserIsNotRoot{},
				Status: check.PASS,
			},
		},
		{
			name: "User is root",
			content: []byte(`FROM ghcr.io/ubuntu:18.04
COPY . /app
USER root
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  LastUserIsNotRoot{},
				Status: check.FAIL,
				Locations: check.Locations{
					load.Range{
						Start: load.Position{Line: 3, Column: 6},
						End:   load.Position{Line: 3, Column: 10},
					},
				},
			},
		},
		{
			name: "User is not root",
			content: []byte(`FROM ghcr.io/ubuntu:18.04
COPY . /app
USER potatoe
RUN make /app
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  LastUserIsNotRoot{},
				Status: check.PASS,
			},
		},
		{
			name: "Last user is root",
			content: []byte(`FROM docker.io/ubuntu:18.04
COPY . /app
USER potatoe
RUN make /app
USER root
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  LastUserIsNotRoot{},
				Status: check.FAIL,
				Locations: check.Locations{
					load.Range{
						Start: load.Position{Line: 5, Column: 6},
						End:   load.Position{Line: 5, Column: 10},
					},
				},
			},
		},
		{
			name: "Last user is not root",
			content: []byte(`FROM ghcr.io/ubuntu:18.04
COPY . /app
USER root
RUN make /app
USER potatoe
CMD python /app/app.py`),
			want: check.CheckResult{
				Check:  LastUserIsNotRoot{},
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

			c := new(LastUserIsNotRoot)
			assert.Equal(t, tt.want, c.Run(*load.NewInput("Dockerfile", nil)))
		})
	}
}
