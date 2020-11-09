// Copyright 2020 Authors of Hubble
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !privileged_tests

package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompileFQDNPattern(t *testing.T) {
	tests := []struct {
		name            string
		fqdnPatterns    []string
		wantErr         bool
		wantErrContains string
		want            string
	}{
		{
			name:            "empty",
			fqdnPatterns:    []string{""},
			wantErr:         true,
			wantErrContains: "empty pattern",
		},
		{
			name:         "simple",
			fqdnPatterns: []string{"cilium.io"},
			want:         `\A(?:cilium\.io)\z`,
		},
		{
			name:         "multiple",
			fqdnPatterns: []string{"cilium.io", "ebpf.io"},
			want:         `\A(?:cilium\.io|ebpf\.io)\z`,
		},
		{
			name:         "star",
			fqdnPatterns: []string{"*.cilium.io"},
			want:         `\A(?:[\-\.0-9a-z]*\.cilium\.io)\z`,
		},
		{
			name:         "trailing_dot",
			fqdnPatterns: []string{"cilium.io."},
			want:         `\A(?:cilium\.io)\z`,
		},
		{
			name:         "spaces",
			fqdnPatterns: []string{"  cilium.io  "},
			want:         `\A(?:cilium\.io)\z`,
		},
		{
			name:         "upper_case",
			fqdnPatterns: []string{"CILIUM.IO"},
			want:         `\A(?:cilium\.io)\z`,
		},
		{
			name:         "spaces_trailing_dot_upper_case",
			fqdnPatterns: []string{"  CILIUM.IO.  "},
			want:         `\A(?:cilium\.io)\z`,
		},
		{
			name:            "empty_after_trim",
			fqdnPatterns:    []string{"  .  "},
			wantErr:         true,
			wantErrContains: "empty pattern",
		},
		{
			name:            "invalid rune",
			fqdnPatterns:    []string{"_"},
			wantErr:         true,
			wantErrContains: "invalid rune in pattern",
		},
		{
			name:            "multiple_trailing_dots",
			fqdnPatterns:    []string{"cilium.io.."},
			wantErr:         true,
			wantErrContains: "multiple trailing dots",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compileFQDNPattern(tt.fqdnPatterns)
			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrContains != "" {
					assert.Contains(t, err.Error(), tt.wantErrContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
		})
	}
}
