package reader

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/google/go-cmp/cmp"
	"testing"
	"web-crawler-go/models"
)

func TestCsvReader_ReadUrlFromSource(t *testing.T) {
	type fields struct {
		readFilePath      string
		headerName        string
		dataValidator     func(string) bool
		processedMappings *hashset.Set
		want              []models.Resource
	}
	type args struct {
		outChan chan *models.Resource
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Invalid file path",
			fields: fields{
				readFilePath:      "test_resources/non_existing_file.csv",
				headerName:        "",
				dataValidator:     nil,
				processedMappings: nil,
				want:              nil,
			},
			args: args{outChan: make(chan *models.Resource, 3)},
		},
		{
			name: "empty csv file",
			fields: fields{
				readFilePath: "./test_resources/empty_csv.csv",
				want:         nil,
			},
			args: args{outChan: make(chan *models.Resource, 3)},
		},
		{
			name: "invalid data",
			fields: fields{
				readFilePath: "./test_resources/invalid_data.csv",
				headerName:   "non_existing_header",
				want:         nil,
			},
			args: args{outChan: make(chan *models.Resource, 3)},
		},
		{
			name: "data validation failed + success",
			args: args{outChan: make(chan *models.Resource, 3)},
			fields: fields{
				readFilePath: "./test_resources/success.csv",
				headerName:   "URL",
				dataValidator: func(s string) bool {
					if s == "amazon.in" {
						return false
					}
					return true
				},
				processedMappings: nil,
				want: []models.Resource{
					{
						ResourceAddress: "google.com",
						ResourceType:    "http",
					},
					{
						ResourceAddress: "flipkart.com",
						ResourceType:    "http",
					},
				},
			},
		},
		{
			name: "success with existing mappings",
			args: args{outChan: make(chan *models.Resource, 3)},
			fields: fields{
				readFilePath:      "./test_resources/success.csv",
				headerName:        "URL",
				processedMappings: hashset.New("amazon.in"),
				want: []models.Resource{
					{
						ResourceAddress: "google.com",
						ResourceType:    "http",
					},
					{
						ResourceAddress: "flipkart.com",
						ResourceType:    "http",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readerInstance := CsvReader{
				readFilePath:      tt.fields.readFilePath,
				headerName:        tt.fields.headerName,
				dataValidator:     tt.fields.dataValidator,
				processedMappings: tt.fields.processedMappings,
			}
			readerInstance.ReadUrlFromSource(tt.args.outChan)
			//close(tt.args.outChan)

			var got []models.Resource

			for resource := range tt.args.outChan {
				got = append(got, *resource)
			}

			if len(got) != len(tt.fields.want) || !cmp.Equal(got, tt.fields.want) {
				t.Errorf("ReadUrlFromSource() got = %v, want %v", got, tt.fields.want)
				return
			}
		})
	}
}
