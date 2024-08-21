package models

import "errors"

var _ Bowl = &PlayerBowl{}
var _ Bowl = &Kalaha{}

type Bowl interface {
	PassBeads(player Player, beads uint8) Player
	GetNext() Bowl
	GetOwner() Player
	GetBeads() uint8
}

type BaseBowl struct {
	Beads   uint8
	Owner   Player
	TheNext Bowl
}

func (pb *BaseBowl) IsEmpty() bool {
	return pb.Beads == 0
}

func (pb *BaseBowl) GetNext() Bowl {
	return pb.TheNext
}

func (pb *BaseBowl) GetOwner() Player {
	return pb.Owner
}

func (pb *BaseBowl) GetBeads() uint8 {
	return pb.Beads
}

type PlayerBowl struct {
	Number uint
	BaseBowl
	Opposite *PlayerBowl
}

func (pb *PlayerBowl) CanMove() bool {
	return pb.Beads > 0
}

func (pb *PlayerBowl) Steal() uint8 {
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

func (pb *PlayerBowl) PassBeads(player Player, beads uint8) Player {
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

func (pb *PlayerBowl) isSteal(player Player, beads uint8) bool {
	return pb.Owner == player && pb.lastBead(beads) && pb.IsEmpty()
}

func (pb *PlayerBowl) lastBead(beads uint8) bool {
	return beads == 1
}

type Kalaha struct {
	Name string
	BaseBowl
}

func (k *Kalaha) PassBeads(player Player, beads uint8) Player {
	if k.Owner == player && beads > 0 {
		k.Beads++
		if beads == 1 {
			return player
		}
		return (k.TheNext).PassBeads(player, beads-1)
	}
	return (k.TheNext).PassBeads(player, beads)

}
