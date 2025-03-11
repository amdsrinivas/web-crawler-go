package writer

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"web-crawler-go/models"
)

func TestFileWriter_ExecuteWrite(t *testing.T) {
	type fields struct {
		TargetDirectory  string
		TruncateMappings bool
	}
	type args struct {
		inChanData []*models.ResourceData
		inChan     chan *models.ResourceData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "check file writer",
			fields: fields{
				TargetDirectory: "resources",
			},
			args: args{
				inChan: make(chan *models.ResourceData, 2),
				inChanData: []*models.ResourceData{
					{
						ResourceAddress: "google.com",
						Data:            []byte("google.com"),
					},
					{
						ResourceAddress: "amazon.com",
						Data:            []byte("amazon.com"),
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.RemoveAll(tt.fields.TargetDirectory)
			writerInstance := &FileWriter{
				TargetDirectory:  tt.fields.TargetDirectory,
				TruncateMappings: tt.fields.TruncateMappings,
			}
			for _, resource := range tt.args.inChanData {
				tt.args.inChan <- resource
			}
			close(tt.args.inChan)
			if got := writerInstance.ExecuteWrite(tt.args.inChan); got != tt.want {
				t.Errorf("ExecuteWrite() = %v, want %v", got, tt.want)
			}

			_, err := os.Open(tt.fields.TargetDirectory)

			if os.IsNotExist(err) {
				t.Errorf("ExecuteWrite() : target directory not created")
				return
			}

			_, err = os.Open(fmt.Sprintf("%s/mapping.csv", tt.fields.TargetDirectory))
			if os.IsNotExist(err) {
				t.Errorf("ExecuteWrite() : target mapping not created")
				return
			}

			files, err := os.ReadDir(tt.fields.TargetDirectory)
			if err != nil || len(files) != 3 {
				t.Errorf("ExecuteWrite() : err: %v, number of files:%d", err, len(files))
				return
			}
		})
	}
}

func Test_getFileWriter(t *testing.T) {
	type args struct {
		opts map[string]any
	}
	tests := []struct {
		name string
		args args
		want Writer
	}{
		{
			name: "Fetch writer",
			args: args{opts: map[string]any{"targetDirectory": "test", "truncateMappings": true}},
			want: &FileWriter{TruncateMappings: true, TargetDirectory: "test"},
		},
		{
			name: "Invalid target directory -- return default",
			args: args{opts: map[string]any{"targetDirectoy": "test", "truncateMappings": false}},
			want: &FileWriter{TruncateMappings: false, TargetDirectory: "resources"},
		},
		{
			name: "Invalid truncate mappings -- return default",
			args: args{opts: map[string]any{"targetDirectory": "test", "truncateMapings": false}},
			want: &FileWriter{TruncateMappings: false, TargetDirectory: "test"},
		},
		{
			name: "Empty opts -- return default",
			args: args{opts: map[string]any{}},
			want: &FileWriter{TruncateMappings: false, TargetDirectory: "resources"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileWriter(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFileWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}
