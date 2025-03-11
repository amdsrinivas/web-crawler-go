package main

import (
	"flag"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/rs/zerolog/log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"web-crawler-go/crawler"
	"web-crawler-go/models"
	"web-crawler-go/reader"
	"web-crawler-go/utils"
	"web-crawler-go/writer"
)

const (
	TargetDirectory = "resources"
)

func main() {

	numCrawlers := flag.Int("num_crawlers", 50, "number of crawlers to run the job with")
	sourcePath := flag.String("source_path", "./urls.csv", "path to the resource list file")
	sourceHeaderName := flag.String("source_header_name", "URL", "Header name for the data column in the source file")
	resumeJob := flag.Bool("resume_job", true, "Resume the job from the last run")
	flag.Parse()

	sigChan := make(chan os.Signal, 1)
	pipelineCompletionChan := make(chan bool, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)

	// Create a buffered channel to control the amount of data loaded into the memory at once.
	readerChan := make(chan *models.Resource, 50)
	writerChan := make(chan *models.ResourceData)

	var mappings *hashset.Set
	if *resumeJob {
		mappings = utils.BuildExistingMappings(fmt.Sprintf("./%s/mapping.csv", TargetDirectory))
	} else {
		mappings = hashset.New()
	}

	// Create pipeline components.

	// Reader : Reads the resources to be crawled from the datasource.
	readerInstance, err := reader.GetReader("csv", map[string]any{
		"mappings":     mappings,
		"readFilepath": *sourcePath,
		"headerName":   *sourceHeaderName,
		"dataValidatorFunc": func(inputUrl string) bool {
			if strings.HasPrefix(inputUrl, "http") {
				_, err := url.ParseRequestURI(inputUrl)
				if err != nil {
					log.Debug().Str("url", inputUrl).Msg("url not valid")
					return false
				}
				return true
			} else {
				_, err := url.ParseRequestURI(fmt.Sprintf("https://%s", inputUrl))
				if err != nil {
					log.Debug().Str("url", inputUrl).Msg("url not valid")
					return false
				}
				return true
			}
		},
	})
	if err != nil {
		log.Fatal().Msg("failed to create the reader")
		os.Exit(1)
	}

	// Crawler : Crawls the resource.
	// Additionally, uses a manager to track the active goroutines and metrics.
	crawlManager := &crawler.CrawlManager{AvailableGoroutines: *numCrawlers}
	crawlerInstance, err := crawler.GetCrawler("http", crawlManager, nil)
	if err != nil {
		log.Fatal().Msg("failed to create the crawler")
		os.Exit(1)
	}

	// Writer : Writes the crawled data to persistence layer.
	writerInstance, err := writer.GetWriter("file", map[string]any{
		"targetDirectory":  TargetDirectory,
		"truncateMappings": !(*resumeJob),
	})
	if err != nil {
		log.Fatal().Msg("failed to create the writer")
		os.Exit(1)
	}

	// Build the pipeline.
	log.Info().Str("source_file", *sourcePath).Msg("starting the crawl job")
	go readerInstance.ReadUrlFromSource(readerChan)
	go crawlerInstance.ExecuteCrawl(readerChan, writerChan)
	go func() {
		pipelineCompletionChan <- writerInstance.ExecuteWrite(writerChan)
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-sigChan:
			crawlManager.ShutdownCrawls()
			log.Info().Msg("shutdown signalled. waiting for writes to complete.")

		case <-pipelineCompletionChan:
			log.Info().Msg("all resources saved. pipeline completed.")
			crawlManager.GenerateReport(true)
			return
		}
	}
}
