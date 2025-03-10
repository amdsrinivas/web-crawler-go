package utils

import (
	"bufio"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"
)

func BuildExistingMappings(mappingFilePath string) *hashset.Set {
	mappings := hashset.New()
	mappingFile, err := os.Open(mappingFilePath)

	if err != nil {
		log.Info().Err(err).Msg("failed to open mapping file. all data will be read.")
		return mappings
	}
	mappingReader := bufio.NewReader(mappingFile)
	currentLine := 0
	for {
		line, err := mappingReader.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		if err == io.EOF {
			break
		}
		currentLine++
		if currentLine > 1 {
			lineData := strings.Split(line, ",")
			if len(lineData) > 1 {
				mappings.Add(lineData[0])
			}
		}
	}
	log.Debug().Any("mappings", mappings).Msg("loaded mappings")
	return mappings
}
