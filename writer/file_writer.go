package writer

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"web-crawler-go/models"

	"github.com/google/uuid"
)

type FileWriter struct {
	TargetDirectory  string
	TruncateMappings bool
}

func getFileWriter(opts map[string]any) Writer {
	targetDirectory, ok := opts["targetDirectory"].(string)
	if !ok {
		log.Warn().Any("opts", opts).Msg("Invalid opts")
		targetDirectory = "resources"
	}

	truncateMappings, ok := opts["truncateMappings"].(bool)
	if !ok {
		log.Warn().Any("opts", opts).Msg("Invalid opts")
	}
	return &FileWriter{
		TargetDirectory:  targetDirectory,
		TruncateMappings: truncateMappings,
	}
}

func (writerInstance *FileWriter) ExecuteWrite(inChan chan *models.ResourceData) bool {
	os.MkdirAll(writerInstance.TargetDirectory, 0777)
	var mappingFileMode int
	if writerInstance.TruncateMappings {
		mappingFileMode = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	} else {
		mappingFileMode = os.O_RDWR | os.O_CREATE | os.O_APPEND
	}
	mappingFile, err := os.OpenFile(fmt.Sprintf("%s/mapping.csv", writerInstance.TargetDirectory), mappingFileMode, 0644)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create mapping file. resumption will not be possible")
		return false
	}
	if mappingFile != nil {
		defer mappingFile.Close()

		stats, _ := mappingFile.Stat()
		if stats.Size() == 0 {
			mappingFile.WriteString("URL,FILE_NAME\n")
		}
	}

	for resource := range inChan {
		resourceIdentifier := uuid.New().String()
		resourceFile, resourceError := os.OpenFile(fmt.Sprintf("%s/%s.txt", writerInstance.TargetDirectory, resourceIdentifier), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if resourceError != nil {
			log.Error().Err(resourceError).Msg("failed to create resource file.")
		}
		if resourceFile != nil {
			resourceFile.Write(resource.Data)
			resourceFile.Close()

			if mappingFile != nil {
				mappingFile.WriteString(fmt.Sprintf("%s,%s\n", resource.ResourceAddress, resourceIdentifier))
			}
		}
	}
	return true
}
