package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"mancala/internal"
	"mancala/pkg"
	"mancala/pkg/models"
	"os"
	"strconv"
	"time"
)

func main() {

	var np int

	app := &cli.App{
		Name:     "Mancala",
		Version:  "0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Manuel Doncel Martos",
				Email: "manueldoncelmartos@gmail.com",
			},
		},
		Usage: "Play Mancala against an AI",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "training",
				Aliases: []string{"t"},
				Usage:   "Number of games to train to the AI",

				Destination: &np,
			},
		},
		EnableBashCompletion: true,
		ArgsUsage:            "-training 100",
		Action: func(cCtx *cli.Context) error {
			n := uint(np)
			if cCtx.Int("training") == 0 {
				n = 500
				logrus.WithField("training", n).Info("Using default number of training games")
			} else {
				n = uint(cCtx.Int("training"))
			}
			logrus.WithField("n", n).Info("Training the AI")
			player2 := internal.Train(n)
			player1 := &models.HumanPlayer{Name: "Player1", BasePlayer: &models.BasePlayer{}}
			player1.Opponent = player2
			player2.Opponent = player1
			ge := pkg.Initialize(player1, player2)

			logrus.Info("Initiating game")
			for ge.PlayerTurn.CanPlay() {
				display(ge)
				var index uint8
				if hp, ok := ge.PlayerTurn.(*models.HumanPlayer); ok {
					for {
						fmt.Printf("%s turn's, select bowl index to move %v: ", hp.Name, ge.PlayerTurn.GetAvailableActions())
						var line string
						_, err := fmt.Scan(&line)
						if err == nil {
							intIndex, err := strconv.Atoi(line)
							if err == nil {
								index = uint8(intIndex)
								break
							}

						}

						fmt.Println("Invalid input. Please enter a valid option.")
					}
				} else {
					state := ge.GetState()
					index = ge.PlayerTurn.(*models.AIPlayer).ChooseAction(state)
					fmt.Printf("AI is playing bowl index: %d\r\n", index)
				}
				if ge.PlayerTurn.CanPlayIndex(index) {
					err := ge.Play(index)
					if err != nil {
						println("error when moving bowl: " + err.Error())
					}
				} else {
					println("Bowl not allowed")
				}

			}
			logrus.WithField("points", ge.GetPoints(player1)).Info("Player 1 points")
			logrus.WithField("points", ge.GetPoints(player2)).Info("Player 2 points")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}

}

func display(engine *pkg.GameEngine) {
	println(" ---- -12- -11- -10-  -9-  -8-  -7- Index ")
	for i := 0; i < 7; i++ {
		var startingBowl models.Bowl = engine.Player2.GetKalaha()
		displayer := pkg.CreateDisplayer(engine.Player1, startingBowl)
		print(displayer.Display(uint(i)))
		startingBowl = startingBowl.GetNext()
		for startingBowl != engine.Player1.GetKalaha().GetNext() {
			displayer := pkg.CreateDisplayer(engine.Player1, startingBowl)
			print(displayer.Display(uint(i)))
			startingBowl = startingBowl.GetNext()
		}
		println()
	}
	println(" Index -0-  -1-  -2-  -3-  -4-  -5-  ---- ")
}
