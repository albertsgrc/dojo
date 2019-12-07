package ai

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

func matchesAnyDescriptors(ai *Ai, descriptors ...Descriptor) bool {
	for _, descriptor := range descriptors {
		if ai.MatchesDescriptor(descriptor) {
			return true
		}
	}

	return false
}

// List :
func List(descriptors ...Descriptor) []*Ai {
	files, err := ioutil.ReadDir(".")

	if err != nil {
		log.Fatal("Could not read current folder contents")
	}

	fileNames := make([]string, len(files))
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	// Get all AIs
	ais := GetAis(fileNames)

	sort.Sort(ByNameAndVersion(ais))

	if len(descriptors) > 0 {
		filteredAis := make([]*Ai, 0)

		for _, ai := range ais {
			if matchesAnyDescriptors(ai, descriptors...) {
				filteredAis = append(filteredAis, ai)
			}
		}

		return filteredAis
	}

	return ais
}

// GetAi ...
func GetAi(descriptor Descriptor) (*Ai, error) {
	ais := List(descriptor)

	if len(ais) == 0 {
		return nil, fmt.Errorf("ai '%s' not found", descriptor.Name)
	}

	return ais[0], nil
}
