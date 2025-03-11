package writer

import (
	"errors"
	"github.com/rs/zerolog/log"
	"web-crawler-go/models"
)

type Writer interface {
	ExecuteWrite(inChan chan *models.ResourceData) bool
}

var writerTypes = map[string]func(map[string]any) Writer{
	"file": getFileWriter,
}

// GetWriter
// Factory function to configure different types of writers.
func GetWriter(writerType string, opts map[string]any) (Writer, error) {
	getWriterFunc, found := writerTypes[writerType]
	if !found {
		log.Error().Str("writer-type", writerType).Msg("unknown writer type")
		return nil, errors.New("unknown writer type")
	}
	return getWriterFunc(opts), nil
}
