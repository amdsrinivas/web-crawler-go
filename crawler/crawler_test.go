package crawler

import (
	"net/http"
	"reflect"
	"testing"
)

func TestGetCrawler(t *testing.T) {
	cm := &CrawlManager{}
	type args struct {
		resourceType string
		manager      *CrawlManager
		client       *http.Client
	}
	tests := []struct {
		name    string
		args    args
		want    Crawler
		wantErr bool
	}{
		{
			name: "check http crawler",
			args: args{
				resourceType: "http",
				manager:      cm,
			},
			want: &HttpCrawler{cm: cm, HttpClient: http.DefaultClient},
		},
		{
			name:    "check unknown crawler",
			args:    args{resourceType: "s3"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrawler(tt.args.resourceType, tt.args.manager, tt.args.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrawler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrawler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
