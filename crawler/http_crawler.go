package crawler

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	"web-crawler-go/models"
)

type HttpCrawler struct {
	HttpClient *http.Client
	cm         *CrawlManager
}

func getHttpCrawler(cm *CrawlManager, client *http.Client) Crawler {
	if client == nil {
		client = http.DefaultClient
	}
	return &HttpCrawler{cm: cm, HttpClient: client}
}

func (crawlerInstance *HttpCrawler) crawl(resource *models.Resource, outChan chan *models.ResourceData) {
	start := time.Now().UnixMilli()

	var requestAddress string
	if strings.HasPrefix(resource.ResourceAddress, "http") {
		requestAddress = resource.ResourceAddress
	} else {
		requestAddress = fmt.Sprintf("http://%s", resource.ResourceAddress)
	}
	resp, err := crawlerInstance.HttpClient.Get(requestAddress)

	if err != nil {
		log.Warn().Err(err).Str("url", requestAddress).Msg("url crawl failed")
		crawlerInstance.cm.RecordFailure()
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Warn().Err(err).Str("url", requestAddress).Msg("url crawl failed")
		crawlerInstance.cm.RecordFailure()
		return
	}
	end := time.Now().UnixMilli()
	crawlerInstance.cm.UpdateRunningAverage(end - start)
	crawlerInstance.cm.RecordSuccess()

	outChan <- &models.ResourceData{
		ResourceAddress: resource.ResourceAddress,
		Data:            body,
	}

}

func (crawlerInstance *HttpCrawler) ExecuteCrawl(inChan chan *models.Resource, outChan chan *models.ResourceData) {
	var writeWaitGroup sync.WaitGroup
	for {
		if crawlerInstance.cm.ReceivedShutdownSignal {
			log.Warn().Msg("shutdown signal received. no more crawls will be spawned.")
			break
		}
		// Registering the goroutine is intentional right at the start to ensure there is no race condition once we
		//check the availability of the goroutines.
		if crawlerInstance.cm.IsGoroutineAvailable() && crawlerInstance.cm.RegisterGoroutine() == nil {
			resource, ok := <-inChan
			if !ok {
				log.Info().Msg("crawl queue is emptied")
				// De-allocate right away as we don't need to process anything.
				crawlerInstance.cm.DeregisterGoroutine()
				break
			}
			writeWaitGroup.Add(1)
			go func() {
				//log.Debug().Msg("crawl goroutine started")
				defer writeWaitGroup.Done()
				defer crawlerInstance.cm.DeregisterGoroutine()
				crawlerInstance.crawl(resource, outChan)
				//log.Debug().Msg("crawl goroutine ended")
			}()
		} else {
			// This needs to be fine-tuned in the production setting.
			// If no goroutine is available, wait for 1 second before trying again.
			time.Sleep(1 * time.Second)
		}
	}

	writeWaitGroup.Wait()
	close(outChan)
}
