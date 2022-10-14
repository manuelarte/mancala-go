package models

type HumanPlayer struct {
	*BasePlayer
	Name string
}

func (hp *HumanPlayer) GetBasePlayer() *BasePlayer {
	return hp.BasePlayer
}
