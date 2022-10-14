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
	Play(index uint8) (Player, error)
	CanPlayIndex(index uint8) bool
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
			if k, ok := b.(*Kalaha); ok && k == bp.GetKalaha() {
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
		if pb, ok := b.(*PlayerBowl); ok && !pb.IsEmpty() && pb.GetOwner().GetBasePlayer() == bp {
			return true
		}
	}
	return false
}

func (bp *BasePlayer) BowlAtIndex(index uint8) (Bowl, error) {
	if index < 0 {
		return nil, errors.New("index not allowed")
	}
	i := uint8(0)
	bg := bp.bowlsGenerator()
	for b := range bg {
		if i == index {
			return b, nil
		}
		i++
	}
	return nil, errors.New("index not allowed")
}

func (bp *BasePlayer) CanPlayIndex(index uint8) bool {
	b, err := bp.BowlAtIndex(index)
	if err != nil {
		return false
	}
	if _, ok := b.(*Kalaha); ok {
		return false
	}
	if pb, ok := b.(*PlayerBowl); ok && pb.Beads <= 0 && pb.GetOwner().GetBasePlayer() == bp {
		return false
	}
	return true
}

func (bp *BasePlayer) Play(index uint8) (Player, error) {
	b, err := bp.BowlAtIndex(index)
	if err != nil {
		return nil, err
	}
	if pb, ok := b.(*PlayerBowl); ok && pb.Owner.GetBasePlayer() == bp {
		return pb.Play()
	}
	return nil, errors.New("index not allowed")
}

func (bp *BasePlayer) GetAvailableActions() []uint8 {
	var availableActions []uint8
	bg := bp.bowlsGenerator()
	index := uint8(0)
	for b := range bg {
		if b.GetBeads() > 0 && b.GetOwner().GetBasePlayer() == bp {
			availableActions = append(availableActions, index)
		}
		index++
	}
	return availableActions
}
