package models

import "errors"

type GameEngine struct {
	Board      Board
	Player1    *Player
	Player2    *Player
	PlayerTurn *Player
}

type Board struct {
	Player1Kalaha Kalaha
	Player2Kalaha Kalaha
	Player1Bowls  *[6]*PlayerBowl
	Player2Bowls  *[6]*PlayerBowl
}

func (ge *GameEngine) Move(index uint) error {
	turn, err := ge.PlayerTurn.Move(index)
	if err != nil {
		return err
	}
	ge.PlayerTurn = turn
	return nil
}

func (ge *GameEngine) Finish() (map[*Player]uint, error) {
	if ge.PlayerTurn.CanPlay() {
		return nil, errors.New("player still allowed to play")
	}
	pointsPerPlayer := map[*Player]uint{ge.Player1: ge.Player1.Kalaha.Beads, ge.Player2: ge.Player2.Kalaha.Beads}
	var b1 Bowl = ge.Player1.StartingBowl
	for b1.TheOwner() != ge.Player2 {
		if pb, ok := b1.(*PlayerBowl); ok {
			pointsPerPlayer[ge.Player1] = pointsPerPlayer[ge.Player1] + pb.Beads
		}
		b1 = b1.Next()
	}
	var b2 Bowl = ge.Player2.StartingBowl
	for b2.TheOwner() != ge.Player1 {
		if pb, ok := b2.(*PlayerBowl); ok {
			pointsPerPlayer[ge.Player2] = pointsPerPlayer[ge.Player2] + pb.Beads
		}
		b2 = b2.Next()
	}
	return pointsPerPlayer, nil
}
