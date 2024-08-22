package models

import (
	"errors"
	"mancala/pkg/utils"
	"math/rand"
)

var _ Player = &HumanPlayer{}
var _ Player = &AIPlayer{}

type Player interface {
	GetBasePlayer() *BasePlayer
	GetStartingBowl() *PlayerBowl
	GetOpponent() Player
	GetKalaha() *Kalaha
	GetAvailableActions() []uint8
	CanPlay() bool
	Play(index uint8) (Player, error)
	CanPlayIndex(index uint8) bool
	SetKalaha(kalaha *Kalaha)
	SetStartingBowl(index *PlayerBowl)
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

func (bp *BasePlayer) SetKalaha(kalaha *Kalaha) {
	bp.Kalaha = kalaha
}

func (bp *BasePlayer) SetStartingBowl(startingBowl *PlayerBowl) {
	bp.StartingBowl = startingBowl
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
	if pb, ok := b.(*PlayerBowl); ok && pb.GetOwner().GetBasePlayer() == bp {
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

type HumanPlayer struct {
	*BasePlayer
	Name string
}

func (hp *HumanPlayer) GetBasePlayer() *BasePlayer {
	return hp.BasePlayer
}

type AIPlayer struct {
	*BasePlayer
	Q       map[utils.Pair]float64
	Alpha   float64
	Epsilon float64
}

func (ai *AIPlayer) GetBasePlayer() *BasePlayer {
	return ai.BasePlayer
}

func (ai *AIPlayer) ChooseAction(state utils.State) uint8 {
	if ai.Epsilon == 0 {
		// Greedy
		return ai.greedy(state)
	} else {
		// Epsilon-Greedy
		epsilonValue := rand.Float64()
		if epsilonValue < ai.Epsilon {
			return ai.getRandomChoice()
		} else {
			return ai.greedy(state)
		}
	}
}

func (ai *AIPlayer) getRandomChoice() uint8 {
	availableActions := ai.GetAvailableActions()
	indexOption := rand.Intn(len(availableActions))
	return availableActions[indexOption]
}

func (ai *AIPlayer) greedy(state utils.State) uint8 {
	availableActions := ai.GetAvailableActions()
	q := map[utils.Pair]float64{}
	for _, a := range availableActions {
		pair := utils.PairFrom(state, a)
		if val, ok := ai.Q[pair]; ok {
			q[pair] = val
		}
	}
	if len(q) == 0 {
		return ai.getRandomChoice()
	} else {
		max := -100.0
		toReturn := uint8(0)
		for key, value := range q {
			if value > max {
				toReturn = key.Action
			}
		}
		return toReturn
	}
}

func (ai *AIPlayer) Update(oldStateAndAction utils.Pair, newState utils.State, reward float64) {
	oldValue := ai.getQValue(oldStateAndAction)
	// best future reward
	bestFuture := ai.bestFutureReward(newState)
	ai.updateQValue(oldStateAndAction, oldValue, reward, bestFuture)

}

func (ai *AIPlayer) getQValue(pair utils.Pair) float64 {
	if val, ok := ai.Q[pair]; ok {
		return val
	}
	return 0.0
}

func (ai *AIPlayer) updateQValue(oldStateAndAction utils.Pair, oldQ float64, reward float64, futureRewards float64) {
	if ai.Q == nil {
		ai.Q = map[utils.Pair]float64{}
	}
	ai.Q[oldStateAndAction] = oldQ + ai.Alpha*(reward+futureRewards-oldQ)
}

func (ai *AIPlayer) bestFutureReward(state utils.State) float64 {
	q := ai.getQForState(state)
	if len(q) == 0 {
		return 0.0
	}
	return getMaxFromQ(q)
}

func (ai *AIPlayer) getQForState(state utils.State) map[utils.Pair]float64 {
	toReturn := map[utils.Pair]float64{}
	for pair, value := range ai.Q {
		if pair.State == state {
			toReturn[pair] = value
		}
	}
	return toReturn
}

func getMaxFromQ(q map[utils.Pair]float64) float64 {
	max := -100.0
	for _, value := range q {
		if value > max {
			max = value
		}
	}
	return max
}
