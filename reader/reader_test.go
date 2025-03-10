package reader

import (
	"reflect"
	"testing"
)

func TestGetReader(t *testing.T) {
	type args struct {
		readerType string
		opts       map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    Reader
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetReader(tt.args.readerType, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetReader() got = %v, want %v", got, tt.want)
			}
		})
	}
}
