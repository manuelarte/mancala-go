package models

import (
	"math/rand"
)

type AIPlayer struct {
	*BasePlayer
	Q       map[Pair]float64
	Alpha   float64
	Epsilon float64
}

func (ai *AIPlayer) GetBasePlayer() *BasePlayer {
	return ai.BasePlayer
}

func (ai *AIPlayer) ChooseAction(state State) uint8 {
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

func (ai *AIPlayer) greedy(state State) uint8 {
	availableActions := ai.GetAvailableActions()
	q := map[Pair]float64{}
	for _, a := range availableActions {
		pair := Pair{
			State:  state,
			Action: a,
		}
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

func (ai *AIPlayer) Update(oldStateAndAction Pair, newState State, reward float64) {
	oldValue := ai.getQValue(oldStateAndAction)
	// best future reward
	bestFuture := ai.bestFutureReward(newState)
	ai.updateQValue(oldStateAndAction, oldValue, reward, bestFuture)

}

func (ai *AIPlayer) getQValue(pair Pair) float64 {
	if val, ok := ai.Q[pair]; ok {
		return val
	}
	return 0.0
}

func (ai *AIPlayer) updateQValue(oldStateAndAction Pair, oldQ float64, reward float64, futureRewards float64) {
	if ai.Q == nil {
		ai.Q = map[Pair]float64{}
	}
	ai.Q[oldStateAndAction] = oldQ + ai.Alpha*(reward+futureRewards-oldQ)
}

func (ai *AIPlayer) bestFutureReward(state State) float64 {
	q := ai.getQForState(state)
	if len(q) == 0 {
		return 0.0
	}
	return getMaxFromQ(q)
}

func (ai *AIPlayer) getQForState(state State) map[Pair]float64 {
	toReturn := map[Pair]float64{}
	for pair, value := range ai.Q {
		if pair.State == state {
			toReturn[pair] = value
		}
	}
	return toReturn
}

func getMaxFromQ(q map[Pair]float64) float64 {
	max := -100.0
	for _, value := range q {
		if value > max {
			max = value
		}
	}
	return max
}
