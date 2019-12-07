package ai

import (
	"strconv"
	"strings"
)

// GetAis ...
func GetAis(fileNames []string) []*Ai {
	ais := make([]*Ai, 0)

	aiMap := make(map[string]*Ai)

	aiFamilyMap := make(map[string]*Family)

	for _, fileName := range fileNames {
		if !IsAiFile(fileName) {
			continue
		}

		fileNameAndExtension := strings.Split(fileName, ".")

		if mapAi, ok := aiMap[fileNameAndExtension[0]]; ok {
			if fileNameAndExtension[1] == "cc" {
				mapAi.FileName = fileName
			}
			continue
		}

		myAi := new(Ai)
		ais = append(ais, myAi)

		ai := ais[len(ais)-1]

		aiMap[fileNameAndExtension[0]] = ai

		ai.FileName = fileName

		trimmedFileName := fileNameAndExtension[0][2:]

		split := strings.Split(trimmedFileName, "_")

		ai.Name = split[0]

		ai.Version = 0
		if len(split) > 1 {
			ai.Version, _ = strconv.Atoi(split[1])
		}

		ai.Description = ""
		if len(split) > 2 {
			ai.Description = split[2]
		}

		if family, ok := aiFamilyMap[ai.Name]; ok {
			family.Add(ai)
		} else {
			family := new(Family)
			family.Add(ai)
			aiFamilyMap[ai.Name] = family
		}
	}

	return ais
}
