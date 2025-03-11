package crawler

import (
	"errors"
	"net/http"
	"web-crawler-go/models"
)

type Crawler interface {
	ExecuteCrawl(readChan chan *models.Resource, outChan chan *models.ResourceData)
}

var crawlerTypes = map[string]func(_ *CrawlManager, _ *http.Client) Crawler{
	"http": getHttpCrawler,
}

func GetCrawler(resourceType string, manager *CrawlManager, client *http.Client) (Crawler, error) {
	getCrawlerFunc, found := crawlerTypes[resourceType]

	if !found {
		return nil, errors.New("unsupported resource type")
	}

	return getCrawlerFunc(manager, client), nil
}
