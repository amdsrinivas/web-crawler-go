package reader

import (
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
		{
			name: "Get CSV reader",
			args: args{readerType: "csv"},
			want: &CsvReader{
				readFilePath:      "./urls.csv",
				headerName:        "URL",
				dataValidator:     nil,
				processedMappings: nil,
			},
			wantErr: false,
		},
		{
			name:    "Get unknown reader",
			args:    args{readerType: "DB"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetReader(tt.args.readerType, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("GetReader() got = %v, want %v", got, tt.want)
			}
		})
	}
}
