package models

import "errors"

var _ Bowl = &PlayerBowl{}
var _ Bowl = &Kalaha{}

type Bowl interface {
	PassBeads(player *Player, beads uint) *Player
	Next() Bowl
}

type PlayerBowl struct {
	Number   uint // just for debugging
	Beads    uint
	Owner    *Player
	TheNext  Bowl
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

func (pb *PlayerBowl) Play() (*Player, error) {
	if pb.Beads == 0 {
		return nil, errors.New("no beads in this bowl")
	}
	previousBeads := pb.Beads
	pb.Beads = 0
	return pb.TheNext.PassBeads(pb.Owner, previousBeads), nil
}

func (pb *PlayerBowl) PassBeads(player *Player, beads uint) *Player {
	if beads == 0 {
		return player.Next
	}
	if pb.isSteal(player, beads) {
		pb.Owner.Kalaha.Beads = pb.Owner.Kalaha.Beads + 1 + pb.Opposite.Steal()
		return player.Next
	}
	pb.Beads++
	return (pb.TheNext).PassBeads(player, beads-1)
}

func (pb *PlayerBowl) isSteal(player *Player, beads uint) bool {
	return pb.Owner == player && pb.lastBead(beads) && pb.IsEmpty()
}

func (pb *PlayerBowl) lastBead(beads uint) bool {
	return beads == 1
}

func (pb *PlayerBowl) IsEmpty() bool {
	return pb.Beads == 0
}

func (pb *PlayerBowl) Next() Bowl {
	return pb.TheNext
}

type Kalaha struct {
	Name    string
	Beads   uint
	Owner   *Player
	TheNext Bowl
}

func (k *Kalaha) PassBeads(player *Player, beads uint) *Player {
	if k.Owner == player && beads > 0 {
		k.Beads++
		if beads == 1 {
			return player
		}
		return (k.TheNext).PassBeads(player, beads-1)
	}
	return (k.TheNext).PassBeads(player, beads)

}

func (k *Kalaha) Next() Bowl {
	return k.TheNext
}
