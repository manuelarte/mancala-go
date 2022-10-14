package models

import (
	"errors"
)

type PlayerBowl struct {
	Number uint
	BaseBowl
	Opposite *PlayerBowl
}

func (pb *PlayerBowl) CanMove() bool {
	return pb.Beads > 0
}

func (pb *PlayerBowl) Steal() uint {
	toReturn := pb.Beads
	pb.Beads = 0
	return toReturn
}

func (pb *PlayerBowl) Play() (Player, error) {
	if pb.Beads == 0 {
		return nil, errors.New("no beads in this bowl")
	}
	previousBeads := pb.Beads
	pb.Beads = 0
	return pb.TheNext.PassBeads(pb.Owner, previousBeads), nil
}

func (pb *PlayerBowl) PassBeads(player Player, beads uint) Player {
	if beads == 0 {
		return player.GetOpponent()
	}
	if pb.isSteal(player, beads) {
		pb.Owner.GetKalaha().Beads = pb.Owner.GetKalaha().Beads + 1 + pb.Opposite.Steal()
		return player.GetOpponent()
	}
	pb.Beads++
	return (pb.TheNext).PassBeads(player, beads-1)
}

func (pb *PlayerBowl) isSteal(player Player, beads uint) bool {
	return pb.Owner == player && pb.lastBead(beads) && pb.IsEmpty()
}

func (pb *PlayerBowl) lastBead(beads uint) bool {
	return beads == 1
}
