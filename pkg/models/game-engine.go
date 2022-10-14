package models

import (
	"errors"
)

type GameEngine struct {
	Board      Board
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

func (ge *GameEngine) Play(index uint) error {
	turn, err := ge.PlayerTurn.Play(index)
	if err != nil {
		return err
	}
	ge.PlayerTurn = turn
	return nil
}

func (ge *GameEngine) Finish() (map[Player]uint, error) {
	if ge.PlayerTurn.CanPlay() {
		return nil, errors.New("player still allowed to play")
	}
	pointsPerPlayer := map[Player]uint{ge.Player1: ge.Player1.GetKalaha().Beads, ge.Player2: ge.Player2.GetKalaha().Beads}
	var b1 Bowl = ge.Player1.GetStartingBowl()
	for b1.GetOwner() != ge.Player2 {
		if pb, ok := b1.(*PlayerBowl); ok {
			pointsPerPlayer[ge.Player1] = pointsPerPlayer[ge.Player1] + pb.Beads
		}
		b1 = b1.GetNext()
	}
	var b2 Bowl = ge.Player2.GetStartingBowl()
	for b2.GetOwner() != ge.Player1 {
		if pb, ok := b2.(*PlayerBowl); ok {
			pointsPerPlayer[ge.Player2] = pointsPerPlayer[ge.Player2] + pb.Beads
		}
		b2 = b2.GetNext()
	}
	return pointsPerPlayer, nil
}
