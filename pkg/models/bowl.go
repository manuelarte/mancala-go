package models

import "errors"

var _ Bowl = &PlayerBowl{}
var _ Bowl = &Kalaha{}

var (
	ErrNoBeadsInThisBowl = errors.New("no beads in this bowl")
)

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

// CanPlay Returns whether the bowl can be played or not
func (pb *PlayerBowl) CanPlay() bool {
	return pb.Beads > 0
}

func (pb *PlayerBowl) Steal() uint8 {
	toReturn := pb.Beads
	pb.Beads = 0
	return toReturn
}

// Play Move the beads of the bowl to the next one
//
// Returns the next player's turn and whether the play was possible or not
func (pb *PlayerBowl) Play() (Player, error) {
	if pb.Beads == 0 {
		return nil, ErrNoBeadsInThisBowl
	}
	previousBeads := pb.Beads
	pb.Beads = 0
	return pb.TheNext.PassBeads(pb.Owner, previousBeads), nil
}

// PassBeads Pass the beads from the bowl to the next one until there are no more beads.
//
// The inputs are the player who started the move and the number of beads that the bowl receives.
// The output is the next player's turn.
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

// Kalaha Struct to represent the Kalaha
type Kalaha struct {
	Name string
	BaseBowl
}

// PassBeads Pass the beads to the next bowl but extracting one if the play was initiated by its owner.
//
// The inputs are the player who initiated the play and the number of beads received by the Kalaha bowl.
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
