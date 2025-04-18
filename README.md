# web-crawler-go

> Go based CLI app to crawl resources from various sources and persist. Currently, supports CSV source, HTTP resources and file based persistence. Built with extensibility to enable addition of different types of sources, resources and persistence layers.

### Stages

#### Read stage
- Reads the records from source into a channel.
- Any type of reader should implement the `Reader` interface.
- Using interface, we can extend the application to different data sources, say `DatabaseReader` without changing the pipeline.

#### Crawl stage
- Processes the read data from the channel and crawls the data.
- Any type of crawler should implement the `Crawler` interface.
- `HttpCrawler` implements the logic to crawl http(s) resources. Can be extended to any type of crawler, say `S3Crawler`.
- `CrawlManager` is used to control the behaviour of crawler goroutines, track the progress and aggregate the metrics.

#### Write stage
- Writes the crawled data to a persistence layer.
- Any type of writer to should implement the `Writer` interface.
- `FileWriter` implements the logic to write the data to file system. Can be extended to any type of writer, say `KafkaWriter`.
- Additionally, writers need to ensure to create the mapping file to keep track of the processed resources between runs.
- `FileWriter` dumps the crawled data to `resources/` directory with a `<uuid>.txt` for each resource along with the mapping file.
- `mappings.csv` file maintains the mapping of the persisted resource and the corresponding resource address.

### Usage
- Ensure Go is installed and GOPATH is set.
```
$go get .
$go run main.go -h
Usage of web-crawler-go:
  -num_crawlers int
        number of crawlers to run the job with (default 50)
  -resume_job
        Resume the job from the last run (default true)
  -source_header_name string
        Header name for the data column in the source file (default "URL")
  -source_path string
        path to the resource list file (default "./urls.csv")
```

-Example: (crawled data will be available in the `resources` directory)
```
go run main.go -num_crawlers=20 -source_header_name=url_host_name -source_path=./dataset/sample_200.csv
```

### Tests
- From the root of the repository:
```
$ go test -cover ./...
```

#### Current coverage
```
        web-crawler-go          coverage: 0.0% of statements
ok      web-crawler-go/crawler  3.100s  coverage: 90.3% of statements
?       web-crawler-go/models   [no test files]
ok      web-crawler-go/reader   0.006s  coverage: 94.8% of statements
ok      web-crawler-go/utils    0.004s  coverage: 100.0% of statements
ok      web-crawler-go/writer   0.005s  coverage: 89.5% of statements
```

### Additional features

#### Shutdown
- Supports graceful shutdown of the system.
- Once an interrupt is triggered, stops spawning new goroutines and waits for the writes to complete and ensures all the crawls in progress are completed and no effort is wasted.

#### State management
- Pipeline uses `mappings.csv` file to maintain the state between runs.
- A shutdown behaves as a pause to the pipeline. Consequent runs picks up from the last stopped stage.
- Pass `-resume_job=false` if you want to restart from scratch. This would clear the mappings file but **will not** remove the crawled data from earlier runs.

#### Report generation
- Generates the report for each run.
- Includes total processed URLs, Error rate (in fraction), Avg. response time (in seconds).

```
#####################################################
Run report:
Total processed URLs:   120
Error rate:     0.3416666666666667
Average response time:  26.121153899059024
#####################################################
```

#### Sample runs
All runs with 50 workers and 50 read buffer. Data is available under `dataset/` directory.
- 999 URLs processed in ~5mins.
```
#####################################################
Run report:
Total processed URLs:   999
Error rate:     0.3103103103103103
Average response time:  18.606190366862805
#####################################################
```

- 199 URLs processed in ~56s
```
#####################################################
Run report:
Total processed URLs:   199
Error rate:     0.19597989949748743
Average response time:  24.841214062872485
#####################################################
```

### Caveats/Assumptions
- `HttpCrawler`'s current implementation does not honor the resource's `robots.txt` or set any standard headers.
- Multiple interrupts are not handled and can lead to unwanted behaviour. Once interrupted, it is better not to interrupt again until it shuts down.
- Logging can be improved to use a custom writer for `zerolog` logger to organize different levels of the log.
- The last entry the input file has to end with `\n`.

### Improvements
- Create a docker image which mounts a volume and updates the file system with the resources.
- `HttpCrawler` should be improved to be built on top of a custom client.