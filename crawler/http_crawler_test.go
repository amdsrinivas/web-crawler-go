package crawler

import (
	"github.com/google/go-cmp/cmp"
	"net/http"
	"testing"
	"web-crawler-go/models"
)

func TestHttpCrawler_ExecuteCrawl(t *testing.T) {
	type fields struct {
		HttpClient *http.Client
		cm         *CrawlManager
	}
	type args struct {
		inChanData []*models.Resource
		inChan     chan *models.Resource
		outChan    chan *models.ResourceData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []models.ResourceData
	}{
		{
			name: "cm with shutdown signal active",
			fields: fields{
				HttpClient: nil,
				cm:         &CrawlManager{ReceivedShutdownSignal: true},
			},
			args: args{
				inChan:  make(chan *models.Resource),
				outChan: make(chan *models.ResourceData),
			},
			want: nil,
		},
		{
			name: "happy path errors + scheme-less URLs and valid URLs",
			fields: fields{
				HttpClient: &http.Client{Transport: MockRoundTripper{}},
				cm:         &CrawlManager{AvailableGoroutines: 1},
			},
			args: args{
				inChanData: []*models.Resource{
					{
						ResourceType:    "http",
						ResourceAddress: "amazon.in",
					},
					{
						ResourceType:    "http",
						ResourceAddress: "google.com",
					},
					{
						ResourceType:    "http",
						ResourceAddress: "http://google.com",
					},
				},
				inChan:  make(chan *models.Resource, 3),
				outChan: make(chan *models.ResourceData, 3),
			},
			want: []models.ResourceData{
				{
					ResourceAddress: "google.com",
					Data:            []byte("google.com"),
				},
				{
					ResourceAddress: "http://google.com",
					Data:            []byte("google.com"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crawlerInstance := &HttpCrawler{
				HttpClient: tt.fields.HttpClient,
				cm:         tt.fields.cm,
			}
			for _, data := range tt.args.inChanData {
				tt.args.inChan <- data
			}
			close(tt.args.inChan)
			crawlerInstance.ExecuteCrawl(tt.args.inChan, tt.args.outChan)
			var got []models.ResourceData

			for resource := range tt.args.outChan {
				got = append(got, *resource)
			}

			if len(got) != len(tt.want) || !cmp.Equal(got, tt.want) {
				t.Errorf("ExecuteCrawl() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
