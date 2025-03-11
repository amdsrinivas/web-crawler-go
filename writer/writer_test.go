package writer

import (
	"reflect"
	"testing"
)

func TestGetWriter(t *testing.T) {
	type args struct {
		writerType string
		opts       map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    Writer
		wantErr bool
	}{
		{
			name: "check file writer",
			args: args{writerType: "file"},
			want: &FileWriter{
				TargetDirectory:  "resources",
				TruncateMappings: false,
			},
			wantErr: false,
		},
		{
			name:    "check unknown writer type",
			args:    args{writerType: "Kafka"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWriter(tt.args.writerType, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWriter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWriter() got = %v, want %v", got, tt.want)
			}
		})
	}
}
