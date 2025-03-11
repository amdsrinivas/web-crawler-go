package reader

import (
	"bufio"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"slices"
	"strings"
	"web-crawler-go/models"
)

func getCsvReader(opts map[string]any) Reader {
	// Default values
	var dataValidatorFunc func(string) bool
	passedFunc, ok := opts["dataValidatorFunc"].(func(string) bool)
	if ok {
		dataValidatorFunc = passedFunc
	} else {
		log.Warn().Msg("dataValidatorFunc is not a function")
	}

	readFilepath, ok := opts["readFilepath"].(string)
	if !ok {
		log.Warn().Any("opts", opts).Msg("Invalid opts")
		readFilepath = "./urls.csv"
	}

	headerName, ok := opts["headerName"].(string)
	if !ok {
		log.Warn().Any("opts", opts).Msg("Invalid opts")
		headerName = "URL"
	}

	mappings, ok := opts["mappings"].(*hashset.Set)
	if !ok {
		log.Warn().Any("opts", opts).Msg("Invalid opts")
	}

	return CsvReader{
		readFilePath:      readFilepath,
		headerName:        headerName,
		dataValidator:     dataValidatorFunc,
		processedMappings: mappings,
	}
}

type CsvReader struct {
	readFilePath      string
	headerName        string
	dataValidator     func(string) bool
	processedMappings *hashset.Set
}

func (readerInstance CsvReader) ReadUrlFromSource(outChan chan *models.Resource) {
	csvFile, err := os.Open(readerInstance.readFilePath)

	if err != nil {
		log.Error().Err(err).Msg("failed to open source file")
		close(outChan)
		return
	}
	defer csvFile.Close()

	fileReader := bufio.NewReader(csvFile)

	var dataIndex, currentLine int
	for {
		line, err := fileReader.ReadString('\n')
		if err == io.EOF {
			if currentLine == 0 {
				log.Debug().Str("file", readerInstance.readFilePath).Msg("empty source file")
			}
			break
		}
		line = strings.Trim(line, "\r\n")
		currentLine = currentLine + 1

		// Process the header line. Assumes only one header line.
		if currentLine == 1 {
			headers := strings.Split(line, ",")
			dataIndex = slices.Index(headers, readerInstance.headerName)
			if dataIndex == -1 {
				log.Warn().Msg("data index not found")
				break
			} else {
				continue
			}
		}

		lineData := strings.Split(line, ",")

		if len(lineData) < dataIndex {
			// Log unable to find the data column.
			log.Warn().Int("rowNo", currentLine).Msg("malformed data")
			continue
		}

		columnData := lineData[dataIndex]
		if readerInstance.processedMappings != nil && readerInstance.processedMappings.Contains(columnData) {
			log.Info().Str("url", columnData).Msg("data already processed. skipping.")
			continue
		}

		if readerInstance.dataValidator != nil && !readerInstance.dataValidator(columnData) {
			// Log malformed data.
			log.Info().Str("url", columnData).Msg("data validation failed. skipping.")
			continue
		}

		// Assumes it is always HTTP resources.
		// Can be extended to support different resources.
		outChan <- &models.Resource{
			ResourceAddress: columnData,
			ResourceType:    "http",
		}

	}
	log.Info().Str("file", readerInstance.readFilePath).Msg("reader finished")
	close(outChan)
}
