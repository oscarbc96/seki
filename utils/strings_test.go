package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindStartAndEndColumn(t *testing.T) {
	type args struct {
		str    string
		substr string
	}
	tests := []struct {
		name         string
		args         args
		wantColStart int
		wantColEnd   int
	}{
		{
			name:         "Substring found",
			args:         args{str: "this a very long string", substr: "very"},
			wantColStart: 8,
			wantColEnd:   11,
		},
		{
			name:         "Substring not found",
			args:         args{str: "this a very long string", substr: "potatoe"},
			wantColStart: -1,
			wantColEnd:   -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colStart, colEnd := FindStartAndEndColumn(tt.args.str, tt.args.substr)
			assert.Equal(t, tt.wantColStart, colStart)
			assert.Equal(t, tt.wantColEnd, colEnd)
		})
	}
}
