package main

import (
	"fmt"
	"mancala/pkg/models"
)

func main() {
	player2 := train(6000)
	player1 := &models.HumanPlayer{Name: "Player1", BasePlayer: &models.BasePlayer{}}
	ge := Initialize(player1, player2)

	println("Initiating game")
	for ge.PlayerTurn.CanPlay() {
		display(ge)
		var index uint8
		if hp, ok := ge.PlayerTurn.(*models.HumanPlayer); ok {
			fmt.Printf("%s turn's, select bowl index to move %v: ", hp.Name, ge.PlayerTurn.GetBasePlayer().GetAvailableActions())
			_, err := fmt.Scanln(&index)
			if err != nil {
				panic("can't parse the input")
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
	println(fmt.Sprintf("%s points: %d", player1.Name, ge.GetPoints(player1)))
	println(fmt.Sprintf("AI points: %d", ge.GetPoints(player2)))

}

func display(engine *models.GameEngine) {
	println(" ---- -12- -11- -10-  -9-  -8-  -7- Index ")
	for i := 0; i < 7; i++ {
		var startingBowl models.Bowl = engine.Player2.GetKalaha()
		displayer := models.CreateDisplayer(engine.Player1, startingBowl)
		print(displayer.Display(uint(i)))
		startingBowl = startingBowl.GetNext()
		for startingBowl != engine.Player1.GetKalaha().GetNext() {
			displayer := models.CreateDisplayer(engine.Player1, startingBowl)
			print(displayer.Display(uint(i)))
			startingBowl = startingBowl.GetNext()
		}
		println()
	}
	println(" Index -0-  -1-  -2-  -3-  -4-  -5-  ---- ")
}

func Initialize(player1 models.Player, player2 models.Player) *models.GameEngine {
	player1.GetBasePlayer().Opponent = player2
	player2.GetBasePlayer().Opponent = player1
	initialBeads := uint8(4)

	player1Kalaha := models.Kalaha{
		Name: "Player1Kalaha",
		BaseBowl: models.BaseBowl{
			Beads: 0,
			Owner: player1,
		},
	}
	player2Kalaha := models.Kalaha{
		Name: "Player2Kalaha",
		BaseBowl: models.BaseBowl{
			Beads: 0,
			Owner: player2,
		},
	}
	player1.GetBasePlayer().Kalaha = &player1Kalaha
	player2.GetBasePlayer().Kalaha = &player2Kalaha

	var player1Bowls [6]*models.PlayerBowl
	var player2Bowls [6]*models.PlayerBowl

	for i := 0; i < 6; i++ {
		var bowl = &models.PlayerBowl{
			Number: uint(i),
			BaseBowl: models.BaseBowl{
				Beads: initialBeads,
				Owner: player1,
			},
		}
		player1Bowls[i] = bowl
		var oppositeBowl = &models.PlayerBowl{
			Number: uint(12 - i),
			BaseBowl: models.BaseBowl{
				Beads: initialBeads,
				Owner: player2,
			},
			Opposite: bowl,
		}
		player2Bowls[5-i] = oppositeBowl
		oppositeBowl.Opposite = bowl
		bowl.Opposite = oppositeBowl
		if i == 5 {
			bowl.TheNext = models.Bowl(&player1Kalaha)
			player1Kalaha.TheNext = oppositeBowl
		}
		if i == 0 {
			oppositeBowl.TheNext = &player2Kalaha
			player2Kalaha.TheNext = bowl
		} else {
			player1Bowls[i-1].TheNext = bowl
			oppositeBowl.TheNext = player2Bowls[6-i]
		}
	}

	startingBowl := player1Bowls[0]
	player1.GetBasePlayer().StartingBowl = startingBowl
	player2.GetBasePlayer().StartingBowl = startingBowl

	return &models.GameEngine{
		Player1:    player1,
		Player2:    player2,
		PlayerTurn: player1,
	}

}

func train(n uint) *models.AIPlayer {
	player1 := &models.AIPlayer{Alpha: 0.5, Epsilon: 0.1, BasePlayer: &models.BasePlayer{}}
	player2 := &models.AIPlayer{Alpha: 0.5, Epsilon: 0.1, BasePlayer: &models.BasePlayer{Opponent: player1}}
	player1.Opponent = player2

	for i := uint(1); i < n; i++ {
		ge := Initialize(player1, player2)
		last := map[models.Player]models.Pair{}
		for ge.PlayerTurn.CanPlay() {
			state := ge.GetState()
			action := ge.PlayerTurn.(*models.AIPlayer).ChooseAction(state)
			last[ge.PlayerTurn] = models.Pair{
				State:  state,
				Action: action,
			}

			err := ge.Play(action)
			if err != nil {
				panic("error in move: " + err.Error())
			}

			if ge.PlayerTurn.CanPlay() {
				if pair, ok := last[ge.PlayerTurn]; ok {
					ge.Player1.(*models.AIPlayer).Update(pair, ge.GetState(), 0)
					ge.Player2.(*models.AIPlayer).Update(pair, ge.GetState(), 0)
				}
			} else {
				// game is finished
				points1 := ge.GetPoints(ge.Player1)
				points2 := ge.GetPoints(ge.Player2)
				var winner models.Player
				var loser models.Player
				if points1 > points2 {
					winner = ge.Player1
					loser = ge.Player2
				} else {
					winner = ge.Player2
					loser = ge.Player1
				}
				ge.Player1.(*models.AIPlayer).Update(last[winner], ge.GetState(), 1)
				ge.Player2.(*models.AIPlayer).Update(last[winner], ge.GetState(), 1)

				ge.Player1.(*models.AIPlayer).Update(last[loser], ge.GetState(), -1)
				ge.Player2.(*models.AIPlayer).Update(last[loser], ge.GetState(), -1)

			}
		}
		if i%100 == 0 {
			fmt.Printf("Training Game %d finished\r\n", i)
		}
	}
	println("train finished")
	return player2

}
