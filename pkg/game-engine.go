package pkg

import (
	"errors"
	"mancala/pkg/models"
	"mancala/pkg/utils"
)

type GameEngine struct {
	Player1    models.Player
	Player2    models.Player
	PlayerTurn models.Player
}

type Board struct {
	Player1Kalaha models.Kalaha
	Player2Kalaha models.Kalaha
	Player1Bowls  *[6]*models.PlayerBowl
	Player2Bowls  *[6]*models.PlayerBowl
}

func (ge *GameEngine) Play(index uint8) error {
	turn, err := ge.PlayerTurn.Play(index)
	if err != nil {
		return err
	}
	ge.PlayerTurn = turn
	return nil
}

func (ge *GameEngine) GetPoints(player models.Player) uint8 {
	if ge.PlayerTurn.CanPlay() {
		return player.GetKalaha().Beads
	} else {
		points := uint8(0)
		var b1 models.Bowl = player.GetStartingBowl()
		if b1.GetOwner() == player {
			points = points + b1.GetBeads()
		}
		b1 = b1.GetNext()
		for b1 != player.GetStartingBowl() {
			if b1.GetOwner() == player {
				points = points + b1.GetBeads()
			}
			b1 = b1.GetNext()
		}
		return points
	}

}

func (ge *GameEngine) Finish() (map[models.Player]uint8, error) {
	if ge.PlayerTurn.CanPlay() {
		return nil, errors.New("player still allowed to play")
	}
	pointsPerPlayer := map[models.Player]uint8{ge.Player1: ge.GetPoints(ge.Player1), ge.Player2: ge.GetPoints(ge.Player2)}
	return pointsPerPlayer, nil
}

func (ge *GameEngine) GetState() utils.State {
	state := [14]uint8{}
	var b models.Bowl = ge.Player1.GetStartingBowl()
	for i := uint8(0); i < 14; i++ {
		state[i] = b.GetBeads()
		b = b.GetNext()
	}
	return state
}

func Initialize(player1 models.Player, player2 models.Player) *GameEngine {
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
	player1.SetKalaha(&player1Kalaha)
	player2.SetKalaha(&player2Kalaha)

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
	player1.SetStartingBowl(startingBowl)
	player2.SetStartingBowl(startingBowl)

	return &GameEngine{
		Player1:    player1,
		Player2:    player2,
		PlayerTurn: player1,
	}

}
