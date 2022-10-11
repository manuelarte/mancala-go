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
	for _, bowl := range ge.Player1.Bowls {
		pointsPerPlayer[ge.Player1] = pointsPerPlayer[ge.Player1] + bowl.Beads
	}
	for _, bowl := range ge.Player2.Bowls {
		pointsPerPlayer[ge.Player2] = pointsPerPlayer[ge.Player2] + bowl.Beads
	}
	return pointsPerPlayer, nil
}
