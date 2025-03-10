package reader

import (
	"errors"
	"web-crawler-go/models"
)

type Reader interface {
	ReadUrlFromSource(outChan chan *models.Resource)
}

var readerTypes = map[string]func(map[string]any) Reader{
	"csv": getCsvReader,
}

func GetReader(readerType string, opts map[string]any) (Reader, error) {

	getReaderFunc, found := readerTypes[readerType]
	if !found {
		return nil, errors.New("unsupported reader type")
	}

	return getReaderFunc(opts), nil
}
