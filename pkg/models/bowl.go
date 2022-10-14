package models

var _ Bowl = &PlayerBowl{}
var _ Bowl = &Kalaha{}

type Bowl interface {
	PassBeads(player Player, beads uint) Player
	GetNext() Bowl
	GetOwner() Player
	GetBeads() uint
}

type BaseBowl struct {
	Beads   uint
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

func (pb *BaseBowl) GetBeads() uint {
	return pb.Beads
}
