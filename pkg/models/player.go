package models

import (
	"errors"
)

var _ Player = &HumanPlayer{}

type Player interface {
	GetBasePlayer() *BasePlayer
	GetStartingBowl() *PlayerBowl
	GetOpponent() Player
	GetKalaha() *Kalaha
	CanPlay() bool
	Play(index uint) (Player, error)
	CanPlayIndex(index uint) bool
}

type BasePlayer struct {
	Kalaha       *Kalaha
	StartingBowl *PlayerBowl
	Opponent     Player
}

func (bp *BasePlayer) GetStartingBowl() *PlayerBowl {
	return bp.StartingBowl
}

func (bp *BasePlayer) GetOpponent() Player {
	return bp.Opponent
}

func (bp *BasePlayer) GetKalaha() *Kalaha {
	return bp.Kalaha
}

type playerBowlChan chan Bowl

func (bp *BasePlayer) bowlsGenerator() playerBowlChan {
	c := make(chan Bowl)
	var b Bowl = bp.StartingBowl
	go func() {
		for {
			if b.GetOwner().GetBasePlayer() != bp {
				close(c)
				return
			}
			c <- b
			b = b.GetNext()
		}
	}()
	return c
}

func (bp *BasePlayer) CanPlay() bool {
	bg := bp.bowlsGenerator()
	for b := range bg {
		if pb, ok := b.(*PlayerBowl); ok && !pb.IsEmpty() {
			return true
		}
	}
	return false
}

func (bp *BasePlayer) BowlAtIndex(index uint) (Bowl, error) {
	if index < 0 {
		return nil, errors.New("index not allowed")
	}
	i := uint(0)
	bg := bp.bowlsGenerator()
	for b := range bg {
		if i == index {
			return b, nil
		}
		i++
	}
	return nil, errors.New("index not allowed")
}

func (bp *BasePlayer) CanPlayIndex(index uint) bool {
	b, err := bp.BowlAtIndex(index)
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

func (bp *BasePlayer) Play(index uint) (Player, error) {
	b, err := bp.BowlAtIndex(index)
	if err != nil {
		return nil, err
	}
	if pb, ok := b.(*PlayerBowl); ok && pb.Owner.GetBasePlayer() == bp {
		return pb.Play()
	}
	return nil, errors.New("index not allowed")
}
