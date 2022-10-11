package models

import "errors"

type Player struct {
	Name   string
	Kalaha *Kalaha
	Bowls  *[6]*PlayerBowl
	Next   *Player
}

func (p *Player) CanPlay() bool {
	for _, bowl := range p.Bowls {
		if bowl.Beads > 0 {
			return true
		}
	}
	return false
}

func (p *Player) CanMove(index uint) bool {
	if p.Bowls[index].Beads <= 0 {
		return false
	}
	if uint(len(p.Bowls)) < index {
		return false
	}
	return true
}

func (p *Player) Move(index uint) (*Player, error) {
	if uint(len(p.Bowls)) < index {
		return p, errors.New("index not allowed")
	}
	return p.Bowls[index].Play()
}
