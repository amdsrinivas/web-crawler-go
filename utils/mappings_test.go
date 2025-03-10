package utils

import (
	"github.com/emirpasic/gods/sets/hashset"
	"reflect"
	"testing"
)

func TestBuildExistingMappings(t *testing.T) {
	type args struct {
		mappingFilePath string
	}
	tests := []struct {
		name string
		args args
		want *hashset.Set
	}{
		{
			name: "empty file",
			args: args{mappingFilePath: "./test_resources/empty_file.csv"},
			want: hashset.New(),
		},
		{
			name: "non-existing file",
			args: args{mappingFilePath: "./test_resources/non_existing.csv"},
			want: hashset.New(),
		},
		{
			name: "file with only header",
			args: args{mappingFilePath: "./test_resources/empty_mappings.csv"},
			want: hashset.New(),
		},
		{
			name: "file with mappings",
			args: args{mappingFilePath: "./test_resources/mappings.csv"},
			want: hashset.New("test_url", "test_url1"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildExistingMappings(tt.args.mappingFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildExistingMappings() = %v, want %v", got, tt.want)
			}
		})
	}
}
