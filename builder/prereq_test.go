package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/suborbital/atmo/directive"
	"github.com/suborbital/subo/project"
)

func TestPrereq_GetCommand(t *testing.T) {
	tests := []struct {
		name    string
		prereq  Prereq
		r       project.RunnableDir
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "successfully expands template",
			prereq: Prereq{
				File:    "_lib/_lib.tar.gz",
				Command: "curl -L https://github.com/suborbital/reactr/archive/v{{ .RunnableDir.Runnable.APIVersion }}.tar.gz -o _lib/_lib.tar.gz",
			},
			r: project.RunnableDir{
				Runnable: &directive.Runnable{
					APIVersion: "0.33.75",
				},
			},
			want:    "curl -L https://github.com/suborbital/reactr/archive/v0.33.75.tar.gz -o _lib/_lib.tar.gz",
			wantErr: assert.NoError,
		},
		{
			name: "errors due to missing data to expand with",
			prereq: Prereq{
				File:    "_lib/_lib.tar.gz",
				Command: "curl -L https://github.com/suborbital/reactr/archive/v{{ .RunnableDir.Runnable.APIVersion }}.tar.gz -o _lib/_lib.tar.gz",
			},
			r: project.RunnableDir{
				Runnable: nil,
			},
			want:    "",
			wantErr: assert.Error,
		},
		{
			name: "successfully expands command with no template tag in it",
			prereq: Prereq{
				File:    "_lib/_lib.tar.gz",
				Command: "curl -L https://github.com/suborbital/reactr/archive/v2.tar.gz -o _lib/_lib.tar.gz",
			},
			r: project.RunnableDir{
				Runnable: &directive.Runnable{
					APIVersion: "0.33.75",
				},
			},
			want:    "curl -L https://github.com/suborbital/reactr/archive/v2.tar.gz -o _lib/_lib.tar.gz",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.prereq.GetCommand(DefaultBuildConfig, tt.r)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
