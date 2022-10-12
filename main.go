package main

import (
	"fmt"
	"mancala/pkg/models"
)

func main() {

	gameEngine := Initialize()
	println("Initiating game")
	for gameEngine.PlayerTurn.CanPlay() == true {
		display(gameEngine)
		print(fmt.Sprintf("%s turn's, select bowl to move [0-5]: ", gameEngine.PlayerTurn.Name))
		var index uint
		_, err := fmt.Scanln(&index)
		if err != nil {
			panic("can't parse the input")
		}
		if gameEngine.PlayerTurn.CanMove(index) {
			err := gameEngine.Move(index)
			if err != nil {
				panic("error when moving bowl")
			}
		} else {
			println("Bowl not allowed")
		}

	}
	points, err := gameEngine.Finish()
	if err != nil {
		panic("can't finish the game")
	}
	println(fmt.Sprintf("points: %v", points))

}

func display(engine *models.GameEngine) {
	for i := 0; i < 7; i++ {
		var startingBowl models.Bowl = engine.Player2.Kalaha
		displayer := models.CreateDisplayer(engine.Player1, startingBowl)
		print(displayer.Display(uint(i)))
		startingBowl = startingBowl.Next()
		for startingBowl != engine.Player1.Kalaha.Next() {
			displayer := models.CreateDisplayer(engine.Player1, startingBowl)
			print(displayer.Display(uint(i)))
			startingBowl = startingBowl.Next()
		}
		println()
	}
}

func Initialize() *models.GameEngine {
	player1 := &models.Player{Name: "Player1"}
	player2 := &models.Player{Name: "Player2", Next: player1}
	player1.Next = player2

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
	player1.Kalaha = &player1Kalaha
	player2.Kalaha = &player2Kalaha

	var player1Bowls [6]*models.PlayerBowl
	var player2Bowls [6]*models.PlayerBowl

	for i := 0; i < 6; i++ {
		var bowl = &models.PlayerBowl{
			Number: uint(i),
			BaseBowl: models.BaseBowl{
				Beads: 4,
				Owner: player1,
			},
		}
		player1Bowls[i] = bowl
		var oppositeBowl = &models.PlayerBowl{
			Number: uint(5 - i),
			BaseBowl: models.BaseBowl{
				Beads: 4,
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

	var board = models.Board{
		Player1Kalaha: player1Kalaha,
		Player2Kalaha: player2Kalaha,
		Player1Bowls:  &player1Bowls,
		Player2Bowls:  &player2Bowls,
	}
	player1.StartingBowl = player1Bowls[0]
	player2.StartingBowl = player2Bowls[0]

	return &models.GameEngine{
		Board:      board,
		Player1:    player1,
		Player2:    player2,
		PlayerTurn: player1,
	}

}
