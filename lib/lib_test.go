/*
Copyright AppsCode Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package lib

import (
	"reflect"
	"testing"
)

func TestSortVersions(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    []string
		wantErr bool
	}{
		{
			name:    "v1 and major",
			args:    []string{"9.6-v1", "9.6"},
			want:    []string{"9.6", "9.6-v1"},
			wantErr: false,
		},
		{
			name:    "v1 and major",
			args:    []string{"9.6-v2", "9.6-v1"},
			want:    []string{"9.6-v1", "9.6-v2"},
			wantErr: false,
		},
		{
			name:    "v1 and major",
			args:    []string{"9.7", "9.6-v1"},
			want:    []string{"9.6-v1", "9.7"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SortVersions(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("SortVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortVersions() got = %v, want %v", got, tt.want)
			}
		})
	}
}
