package models

import "errors"

type Player struct {
	Name         string
	Kalaha       *Kalaha
	StartingBowl *PlayerBowl
	Next         *Player
}

func (p *Player) CanPlay() bool {
	var b Bowl = p.StartingBowl
	for b.TheOwner() == p {
		b = b.Next()
		if k, ok := b.(*Kalaha); ok {
			b = k.Next()
		}
		if pb, ok := b.(*PlayerBowl); ok && !pb.IsEmpty() {
			return true
		}
	}
	return false
}

func (p *Player) BowlAtIndex(index uint) (Bowl, error) {
	if index < 0 {
		return nil, errors.New("index not allowed")
	}
	i := index
	var b Bowl = p.StartingBowl
	for i > 0 {
		b = b.Next()
		i--
	}
	return b, nil
}

func (p *Player) CanMove(index uint) bool {
	b, err := p.BowlAtIndex(index)
	if err != nil {
		return false
	}
	if _, ok := b.(*Kalaha); ok {
		return false
	}
	if pb, ok := b.(*PlayerBowl); ok && pb.Beads <= 0 {
		return false
	}
	return true
}

func (p *Player) Move(index uint) (*Player, error) {
	b, err := p.BowlAtIndex(index)
	if err != nil {
		return nil, err
	}
	if pb, ok := b.(*PlayerBowl); ok && pb.Owner == p {
		return pb.Play()
	}
	return p, errors.New("index not allowed")
}
