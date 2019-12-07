package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sort"
	"time"

	"github.com/korovkin/limiter"

	"github.com/albertsgrc/dojo/v2/ai"
)

type aiResults struct {
	NumGamesAtPlaceOrBetter []int
	NumGames                int
	Scores                  []int
	NumWinsEvaluated        int
	Elo                     int
}

type EvaluationResult struct {
	Player                  string
	NumGamesAtPlaceOrBetter []int
	Scores                  []int
	NumWinsEvaluated        int
	Elo                     int
}

type gameResultError struct {
	result GameResult
	err    error
}

// ByEloDescending ...
type ByEloDescending []*EvaluationResult

func (a ByEloDescending) Len() int {
	return len(a)
}

func (a ByEloDescending) Less(i, j int) bool {
	return a[i].Elo > a[j].Elo
}

func (a ByEloDescending) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func winProbability(elo1, elo2 int) float64 {
	return 1.0 / (1 + math.Pow(10, float64(elo1-elo2)/400.0))
}

func updateElos(winner, loser *aiResults) {
	pLoser := winProbability(winner.Elo, loser.Elo)
	pWinner := winProbability(loser.Elo, winner.Elo)

	winner.Elo += int(11 * (1 - pWinner))
	loser.Elo += int(11 * (0 - pLoser))
}

func runGame(randGenTime *rand.Rand, ais []*ai.Ai, limit *limiter.ConcurrencyLimiter, gameResults chan gameResultError) {
	descriptors := []string{}

	playerSet := make(map[string]bool)

	for player := 0; player < 4; player++ {
		descriptor := ais[randGenTime.Intn(len(ais))].Descriptor()
		if player < len(ais) {
			_, ok := playerSet[descriptor]
			for ok {
				descriptor = ais[randGenTime.Intn(len(ais))].Descriptor()
				_, ok = playerSet[descriptor]
			}
		}

		playerSet[descriptor] = true
		descriptors = append(descriptors, descriptor)

	}

	limit.Execute(func() {
		gameResult, err := Run(descriptors, "time", true, false)
		gameResults <- gameResultError{gameResult, err}
	})
}

func processResult(evaluatedAi *ai.Ai, res gameResultError, aiToResults map[string]*aiResults, errChan chan error, onGameFinished func()) {
	if res.err != nil {
		errChan <- res.err
	}

	gameResult := res.result

	scores := make(map[string]int)
	for i, player := range gameResult.Players {
		if gameResult.Scores[i] >= scores[player] {
			scores[player] = gameResult.Scores[i]
		}
	}

	for player, score := range scores {
		var evaluations *aiResults
		var ok bool
		if evaluations, ok = aiToResults[player]; !ok {
			evaluations = new(aiResults)
			evaluations.Scores = make([]int, 0)
			evaluations.NumGamesAtPlaceOrBetter = make([]int, 3)
			evaluations.Elo = 1500
		}

		evaluations.Scores = append(evaluations.Scores, score)
		aiToResults[player] = evaluations
	}

	winner := gameResult.Players[gameResult.Winner]

	if winner == evaluatedAi.PlayerName() {
		for player := range scores {
			aiToResults[player].NumWinsEvaluated++
		}
	}

	players := gameResult.PlayersSorted
	p0 := aiToResults[players[0]]
	p1 := aiToResults[players[1]]
	p2 := aiToResults[players[2]]
	p3 := aiToResults[players[3]]

	updateElos(p0, p1)
	updateElos(p0, p2)
	updateElos(p0, p3)
	updateElos(p1, p2)
	updateElos(p1, p3)
	updateElos(p2, p3)

	for i, player := range gameResult.PlayersSorted {
		for j := i; j < 3; j++ {
			aiToResults[player].NumGamesAtPlaceOrBetter[j]++
		}
	}

	onGameFinished()
}

// Evaluate ...
func Evaluate(evaluatedAi *ai.Ai, numGames int, againstDescriptors []string, onGameFinished func()) ([]*EvaluationResult, error) {
	numDescriptors := len(againstDescriptors)

	if numDescriptors == 0 {
		return nil, fmt.Errorf("evaluate received no against descriptors")
	}

	againstDescriptorsValue := make([]ai.Descriptor, len(againstDescriptors)+1)
	againstDescriptorsValue[0] = ai.DescriptorFromString(evaluatedAi.Descriptor())
	for i, descriptor := range againstDescriptors {
		againstDescriptorsValue[i+1] = ai.DescriptorFromString(descriptor)
	}

	ais := ai.List(againstDescriptorsValue...)

	aiToResults := make(map[string]*aiResults)
	gameResults := make(chan gameResultError, 200)
	errChan := make(chan error)

	go func() {
		for res := range gameResults {
			processResult(evaluatedAi, res, aiToResults, errChan, onGameFinished)
		}

		close(errChan)
	}()

	limit := limiter.NewConcurrencyLimiter(runtime.NumCPU())
	s := rand.NewSource(time.Now().UnixNano() / 1000)
	randGenTime := rand.New(s)

	for game := 0; game < numGames; game++ {
		runGame(randGenTime, ais, limit, gameResults)
	}

	limit.Wait()
	close(gameResults)

	for err := range errChan {
		return nil, err
	}

	evaluationResults := make([]*EvaluationResult, 0)
	for player, aiResults := range aiToResults {

		evaluationResult := new(EvaluationResult)
		evaluationResult.Player = player
		evaluationResult.NumGamesAtPlaceOrBetter = aiResults.NumGamesAtPlaceOrBetter
		evaluationResult.Scores = aiResults.Scores
		evaluationResult.NumWinsEvaluated = aiResults.NumWinsEvaluated
		evaluationResult.Elo = aiResults.Elo
		evaluationResults = append(evaluationResults, evaluationResult)
	}

	sort.Sort(ByEloDescending(evaluationResults))

	return evaluationResults, nil
}
