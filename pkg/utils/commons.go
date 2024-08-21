package utils

type Pair struct {
	State  State
	Action uint8 // between 0-13
}

type State [14]uint8
