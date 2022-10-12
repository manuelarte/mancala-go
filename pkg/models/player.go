package models

import "errors"

type Player struct {
	Name         string
	Kalaha       *Kalaha
	StartingBowl *PlayerBowl
	Next         *Player
}

func (p *Player) CanPlay() bool {
	bg := p.BowlsGenerator()
	for b := range bg {
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
	i := uint(0)
	bg := p.BowlsGenerator()
	for b := range bg {
		if i == index {
			return b, nil
		}
		i++
	}
	return nil, errors.New("index not allowed")
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

type bowlChan chan Bowl

func (p *Player) BowlsGenerator() bowlChan {
	c := make(chan Bowl)
	var b Bowl = p.StartingBowl
	go func() {
		for {
			if b.TheOwner() != p {
				close(c)
				return
			}
			c <- b
			b = b.Next()
		}
	}()
	return c
}

func (b bowlChan) Next() Bowl {
	c, ok := <-b
	if !ok {
		return nil
	}
	return c
}
