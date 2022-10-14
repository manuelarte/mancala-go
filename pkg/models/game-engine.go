package models

import (
	"errors"
)

type GameEngine struct {
	Player1    Player
	Player2    Player
	PlayerTurn Player
}

type Board struct {
	Player1Kalaha Kalaha
	Player2Kalaha Kalaha
	Player1Bowls  *[6]*PlayerBowl
	Player2Bowls  *[6]*PlayerBowl
}

func (ge *GameEngine) Play(index uint8) error {
	turn, err := ge.PlayerTurn.Play(index)
	if err != nil {
		return err
	}
	ge.PlayerTurn = turn
	return nil
}

func (ge *GameEngine) GetPoints(player Player) uint8 {
	if ge.PlayerTurn.CanPlay() {
		return player.GetKalaha().Beads
	} else {
		points := uint8(0)
		var b1 Bowl = player.GetStartingBowl()
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

func (ge *GameEngine) Finish() (map[Player]uint8, error) {
	if ge.PlayerTurn.CanPlay() {
		return nil, errors.New("player still allowed to play")
	}
	pointsPerPlayer := map[Player]uint8{ge.Player1: ge.GetPoints(ge.Player1), ge.Player2: ge.GetPoints(ge.Player2)}
	return pointsPerPlayer, nil
}

func (ge *GameEngine) GetState() State {
	state := [14]uint8{}
	var b Bowl = ge.Player1.GetBasePlayer().GetStartingBowl()
	for i := uint8(0); i < 14; i++ {
		state[i] = b.GetBeads()
		b = b.GetNext()
	}
	return state
}
