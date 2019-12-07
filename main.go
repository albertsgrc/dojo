package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/progress"
	"github.com/jedib0t/go-pretty/table"
	"github.com/montanaflynn/stats"

	plist "github.com/jedib0t/go-pretty/list"
	"github.com/jedib0t/go-pretty/text"

	"github.com/albertsgrc/dojo/v2/ai"
	"github.com/albertsgrc/dojo/v2/utils"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func list(c *cli.Context) error {
	descriptors := make([]ai.Descriptor, c.NArg())

	for i := 0; i < c.NArg(); i++ {
		arg := c.Args().Get(i)
		descriptors[i] = ai.DescriptorFromString(arg)
	}

	ais := ai.List(descriptors...)

	if len(ais) == 0 {
		fmt.Println("No AIs found")
	}

	l := plist.NewWriter()
	l.SetStyle(plist.StyleConnectedRounded)
	previousName := ""
	for _, x := range ais {
		if previousName != x.Name {
			l.UnIndent()
			l.AppendItem(text.Colors{text.Bold, text.BgBlack, text.FgRed}.Sprint(x.Family.Name))
			l.Indent()
		}

		isSelected := ""
		if x.MatchesDescriptor(ai.DescriptorFromString(c.String("ai"))) {
			isSelected = " ✨"
		}

		item := x.FileName + isSelected

		if x.Version == x.Family.LastVersion.Version {
			item = text.Bold.Sprint(item)
		}

		l.AppendItem(item)

		previousName = x.Name
	}
	fmt.Println(l.Render())

	return nil
}

func newVersion(c *cli.Context) error {
	from := c.String("from")
	if len(from) == 0 {
		from = c.String("ai")
	}

	myAi, err := ai.GetAi(ai.DescriptorFromString(from))

	if err != nil {
		return err
	}

	description := ""
	if c.NArg() > 0 {
		description = c.Args().First()
	}

	baseAi, version := ai.NewVersion(myAi, description)

	fmt.Printf(
		"🚀 created version %s for AI %s based on %s\n",
		text.Bold.Sprint(version),
		text.Bold.Sprint(baseAi.Name),
		text.Bold.Sprint(baseAi.FileName))

	return nil
}

func run(c *cli.Context) error {
	pw := progress.NewWriter()
	pw.SetTrackerLength(25)
	pw.ShowOverallTracker(false)
	pw.ShowTime(true)
	pw.ShowTracker(false)
	pw.ShowValue(false)
	pw.SetMessageWidth(24)
	pw.SetNumTrackersExpected(1)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 100)
	pw.Style().Colors = progress.StyleColorsExample

	if !c.Bool("print-output") {
		go pw.Render()
	}

	fmt.Printf("Compiling ... ")
	err := utils.Compile()
	fmt.Printf("done\n")

	if err != nil {
		return err
	}

	trackerRun := progress.Tracker{Message: "Running game"}
	pw.AppendTracker(&trackerRun)
	gameResult, errRun := Run(c.StringSlice("players"), c.String("seed"), c.Bool("shuffle"), c.Bool("print-output"))
	trackerRun.MarkAsDone()

	pw.Stop()

	if errRun != nil {
		return errRun
	}

	fmt.Println()
	fmt.Print(gameResult)

	return nil
}

func ff(x float64) string {
	return fmt.Sprintf("%.2f", x)
}

func fr(x interface{}, isSpecial bool) string {
	if isSpecial {
		return text.Bold.Sprint(x)
	} else {
		return fmt.Sprint(x)
	}
}

func fw(x float64) string {
	var styler text.Color
	if x <= 25 {
		styler = text.FgRed
	} else if x <= 50 {
		styler = text.FgYellow
	} else if x <= 75 {
		styler = text.FgGreen
	} else {
		styler = text.FgBlue
	}

	return styler.Sprintf("%.2f", x)
}

func evaluate(c *cli.Context) error {
	myAi, err := ai.GetAi(ai.DescriptorFromString(c.String("ai")))

	if err != nil {
		return err
	}

	pw := progress.NewWriter()
	pw.SetTrackerLength(20)
	//pw.ShowOverallTracker(true)
	pw.ShowOverallTracker(true)
	pw.ShowTime(false)
	pw.ShowTracker(false)
	pw.ShowValue(true)
	pw.SetMessageWidth(18)
	pw.SetNumTrackersExpected(1)
	pw.SetStyle(progress.StyleDefault)
	pw.SetTrackerPosition(progress.PositionRight)
	pw.SetUpdateFrequency(time.Millisecond * 1000)
	pw.SetAutoStop(true)
	pw.Style().Colors = progress.StyleColorsExample
	pw.Style().Chars = progress.StyleCharsCircle

	go pw.Render()

	fmt.Printf("Compiling ... ")

	//pw.AppendTracker(&trackerCompile)
	err = utils.Compile()
	fmt.Printf("done\n")

	if err != nil {
		return err
	}

	numGames := c.Int("games")
	trackerMessage := fmt.Sprintf("Running %d games", numGames)
	trackerEvaluate := progress.Tracker{Message: trackerMessage, Total: int64(numGames)}
	pw.AppendTracker(&trackerEvaluate)

	evaluations, err := Evaluate(myAi, numGames, c.StringSlice("against"), func() {
		trackerEvaluate.Increment(1)
	})
	trackerEvaluate.MarkAsDone()
	pw.Stop()

	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredGreenWhiteOnBlack)
	t.SetTitle("Ranking")
	t.AppendHeader(table.Row{"#", "AI", "Elo", "1st%", "<=2nd%", "<=3rd%", "EvWin%", "Score", "95%", "99%", "Games"})

	for i, evaluation := range evaluations {
		numGames := len(evaluation.Scores)

		data := make([]float64, numGames)

		for i, score := range evaluation.Scores {
			data[i] = float64(score)
		}

		firstPercentage := 100 * float64(evaluation.NumGamesAtPlaceOrBetter[0]) / float64(numGames)
		secondPercentage := 100 * float64(evaluation.NumGamesAtPlaceOrBetter[1]) / float64(numGames)
		thirdPercentage := 100 * float64(evaluation.NumGamesAtPlaceOrBetter[2]) / float64(numGames)
		avgScore, _ := stats.Mean(data)
		stdevScore, _ := stats.StandardDeviation(data)
		percentile95, _ := stats.Percentile(data, 95)
		percentile99, _ := stats.Percentile(data, 99)
		winPercentageEvaluated := 100 * float64(evaluation.NumWinsEvaluated) / float64(numGames)

		playerPrefix := "  "

		isSpecial := evaluation.Player == myAi.PlayerName() || i == 0

		if evaluation.Player == myAi.PlayerName() {
			playerPrefix = "✨"
			winPercentageEvaluated = 0
		}

		t.AppendRow(table.Row{
			i + 1,
			fr(playerPrefix+evaluation.Player, isSpecial),
			fr(evaluation.Elo, isSpecial),
			fr(fw(firstPercentage), isSpecial),
			fr(fw(secondPercentage), isSpecial),
			fr(fw(thirdPercentage), isSpecial),
			fr(ff(winPercentageEvaluated), isSpecial),
			fr(fmt.Sprintf(`%.2f ± %.2f%s`, avgScore, 100*stdevScore/avgScore, "%"), isSpecial),
			fr(ff(percentile95), isSpecial),
			fr(ff(percentile99), isSpecial),
			fr(numGames, isSpecial),
		})
	}

	t.Render()

	return nil
}

func before(c *cli.Context) error {
	return altsrc.InitInputSourceWithContext(c.Command.Flags, altsrc.NewTomlSourceFromFlagFunc("config"))(c)
}

func main() {
	if _, err := os.Stat("dojo.toml"); os.IsNotExist(err) {
		fmt.Println("Config file dojo.toml not found in current directory. See https://github.com/albertsgrc/dojo/dojo.toml for an example")
		os.Exit(1)
	}

	generalFlags := []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "current-ai",
			Aliases:     []string{"ai"},
			Usage:       "the AI descriptor of the AI you are currently working with",
			DefaultText: "Demo",
			Value:       "Demo",
		}),
		&cli.StringFlag{
			Name:        "config",
			Usage:       "path to the configuration file",
			DefaultText: "dojo.toml file in current directory",
			Value:       "dojo.toml",
		},
	}

	commands := []*cli.Command{
		{
			Name:  "ai",
			Usage: "AI related operations",
			Subcommands: []*cli.Command{
				{
					Name:   "list",
					Usage:  "list AIs",
					Action: list,
				},
				{
					Name:  "new",
					Usage: "create a new version of an ai",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "from",
							Usage:       "the new AI's source code will be copied from `AI_FROM`",
							DefaultText: "current ai",
						},
					},
					Before: before,
					Action: newVersion,
				},
			},
		},
		{
			Name:  "run",
			Usage: "run a game with given AIs",
			Flags: []cli.Flag{
				altsrc.NewBoolFlag(&cli.BoolFlag{
					Name:        "run.shuffle",
					Aliases:     []string{"shuffle"},
					Usage:       "shuffle the positions of the AIs in the map",
					DefaultText: "false",
					Value:       false,
				}),
				altsrc.NewBoolFlag(&cli.BoolFlag{
					Name:        "run.print-output",
					Aliases:     []string{"print-output"},
					Usage:       "print the output from the run command",
					DefaultText: "true",
					Value:       false,
				}),
				altsrc.NewStringFlag(&cli.StringFlag{
					Name:        "run.seed",
					Aliases:     []string{"seed"},
					Usage:       "set the random seed, either a number or the string 'time'",
					DefaultText: "0",
					Value:       "0",
				}),
				altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
					Name:        "run.players",
					Aliases:     []string{"players", "p"},
					Usage:       "set the game payers, e.g. -p Demo --p Dummy -p Dummy -p Dummy",
					DefaultText: "<CurrentAi> Dummy Dummy Dummy",
				}),
			},
			Before: before,
			Action: run,
		},
		{
			Name:   "evaluate",
			Usage:  "evaluate an AI's performance by playing against other AIs",
			Before: before,
			Flags: []cli.Flag{
				altsrc.NewIntFlag(&cli.IntFlag{
					Name:        "evaluate.games",
					Aliases:     []string{"games"},
					Usage:       "number of games that the ai will be evaluated on",
					DefaultText: "200",
					Value:       200,
				}),
				altsrc.NewStringSliceFlag(&cli.StringSliceFlag{
					Name:        "evaluate.against",
					Aliases:     []string{"against"},
					Usage:       "opponent AIs will be chosen from the pool described by `AI_DESCR`",
					DefaultText: ":",
					Value:       cli.NewStringSlice(":"),
				}),
			},
			Action: evaluate,
		},
	}

	app := &cli.App{
		Name:     "dojo",
		HelpName: "dojo",
		Usage:    "manage versioning, running and evaluating your EDA game AIs",
		Before: altsrc.InitInputSourceWithContext(
			generalFlags,
			altsrc.NewTomlSourceFromFlagFunc("config"),
		),
		Flags:    generalFlags,
		Commands: commands,
	}

	err := app.Run(os.Args)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
