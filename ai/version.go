package ai

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// NewVersion ...
func NewVersion(ai *Ai, description string) (*Ai, int) {
	if !strings.HasSuffix(ai.FileName, ".cc") {
		fmt.Println("The AI ", ai.Name, " does not have a source file")
		os.Exit(1)
	}

	newVersion := ai.Family.LastVersion.Version + 1

	playerName := fmt.Sprintf("%s_%d", ai.Name, newVersion)

	filePlayerName := playerName
	if len(description) > 0 {
		filePlayerName = filePlayerName + "_" + description
	}

	fileName := fmt.Sprintf("AI%s.cc", filePlayerName)

	content, _ := ioutil.ReadFile(ai.FileName)

	re := regexp.MustCompile(`#define PLAYER_NAME [\w\d_]+`)

	newContent := re.ReplaceAllString(
		string(content),
		fmt.Sprintf("#define PLAYER_NAME %s", playerName))

	ioutil.WriteFile(fileName, []byte(newContent), 0644)

	return ai, newVersion
}
