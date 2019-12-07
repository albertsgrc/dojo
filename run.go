package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/albertsgrc/dojo/v2/utils"
	"github.com/jedib0t/go-pretty/text"

	"github.com/albertsgrc/dojo/v2/ai"
)

// GameResult ...
type GameResult struct {
	Players       []string
	PlayersSorted []string
	Scores        []int
	Winner        int
}

type ByScoreDescending GameResult

func (a ByScoreDescending) Len() int {
	return len(a.Players)
}

func (a ByScoreDescending) Less(i, j int) bool {
	return a.Scores[i] > a.Scores[j]
}

func (a ByScoreDescending) Swap(i, j int) {
	a.PlayersSorted[i], a.PlayersSorted[j] = a.PlayersSorted[j], a.PlayersSorted[i]
}

func (gr GameResult) String() string {
	s := ""

	for i := 0; i < 4; i++ {
		nameStyler := text.Colors{}
		var valueStyler text.Colors

		prefix := "   "

		score := gr.Scores[i]

		if score < 300 {
			valueStyler = text.Colors{text.FgRed, text.BgBlack}
		} else if score < 600 {
			valueStyler = text.Colors{text.FgYellow, text.BgBlack}
		} else if score < 1000 {
			valueStyler = text.Colors{text.FgGreen, text.BgBlack}
		} else {
			valueStyler = text.Colors{text.FgBlue, text.BgBlack}
		}

		if i == gr.Winner {
			nameStyler = text.Colors{text.Bold}
			prefix = "✌️  "
		}

		s += fmt.Sprintln(
			nameStyler.Sprint(prefix, text.AlignLeft.Apply(gr.Players[i], 14),
				valueStyler.Sprint(score)))
	}

	return s
}

func parseGameResult(output string) GameResult {
	lines := strings.Split(output, "\n")
	lines = lines[len(lines)-7 : len(lines)-3]

	r, _ := regexp.Compile(`player ([_\d\w]+) got score (\d+)`)
	gameResult := GameResult{
		Players: make([]string, 4),
		Scores:  make([]int, 4),
	}

	maxScore := 0
	for i, line := range lines {
		res := r.FindAllStringSubmatch(line, 2)

		gameResult.Players[i] = res[0][1]
		gameResult.Scores[i], _ = strconv.Atoi(res[0][2])

		if gameResult.Scores[i] > maxScore {
			maxScore = gameResult.Scores[i]
			gameResult.Winner = i
		}
	}

	gameResult.PlayersSorted = make([]string, 4)
	copy(gameResult.PlayersSorted, gameResult.Players)

	sort.Sort(ByScoreDescending(gameResult))

	return gameResult
}

// Run ...
func Run(playerDescriptors []string, seed string, shuffle bool, printOutput bool) (GameResult, error) {
	if len(playerDescriptors) != 4 {
		return GameResult{}, fmt.Errorf("Invalid number of players '%d'", len(playerDescriptors))
	}

	s := rand.NewSource(time.Now().UnixNano() / 1000)
	randGenTime := rand.New(s)

	players := make([]string, 4)

	for i, player := range playerDescriptors {
		ais := ai.List(ai.DescriptorFromString(player))

		if len(ais) == 0 {
			return GameResult{}, fmt.Errorf("No AIs found for player %d with descriptor %s", i, player)
		}

		ai := ais[randGenTime.Intn(len(ais))]

		players[i] = ai.PlayerName()
	}

	if shuffle {
		randGenTime.Shuffle(len(players), func(i, j int) {
			players[i], players[j] = players[j], players[i]
		})
	}

	if seed == "time" {
		seed = strconv.FormatInt((time.Now().UnixNano()/1000)%2147479307, 10)
	}

	_, stderr, err := utils.Exec("Game", printOutput,
		players[0], players[1], players[2], players[3],
		"-s", seed, "-i", "default.cnf", "-o", "default.res")

	if err != nil {
		fmt.Println(stderr)
		utils.Error("Running the game failed, see error above ^")
		os.Exit(1)

		return GameResult{}, err
	}

	gameResult := parseGameResult(stderr)

	return gameResult, nil
}
