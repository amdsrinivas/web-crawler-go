package crawler

import (
	"errors"
	"web-crawler-go/models"
)

type Crawler interface {
	ExecuteCrawl(readChan chan *models.Resource, outChan chan *models.ResourceData)
}

var crawlerTypes = map[string]func(manager *CrawlManager) Crawler{
	"http": getHttpCrawler,
}

func GetCrawler(resourceType string, manager *CrawlManager) (Crawler, error) {
	getCrawlerFunc, found := crawlerTypes[resourceType]

	if !found {
		return nil, errors.New("unsupported resource type")
	}

	return getCrawlerFunc(manager), nil
}
