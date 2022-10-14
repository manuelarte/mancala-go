package models

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
