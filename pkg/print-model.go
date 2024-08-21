package pkg

import (
	"fmt"
	"mancala/pkg/models"
)

var _ Displayer = &PrintPlayerBowl{}
var _ Displayer = &PrintKalaha{}

type Displayer interface {
	Display(line uint) string
}

type PrintPlayerBowl struct {
	*models.PlayerBowl
	Player1 models.Player
}

func (ppb *PrintPlayerBowl) Display(line uint) string {
	lineSeparator := " --- "
	beadsPlayer1 := fmt.Sprintf("| %d |", ppb.Beads)
	beadsPlayer2 := fmt.Sprintf("| %d |", ppb.Opposite.Beads)
	return [7]string{lineSeparator, beadsPlayer2, lineSeparator, "-----", lineSeparator, beadsPlayer1, lineSeparator}[line]
}

type PrintKalaha struct {
	*models.Kalaha
	Player1 models.Player
}

func (pk *PrintKalaha) Display(line uint) string {
	string0 := " ---- "
	stringEmpty := "|    |"
	var string5 string
	if pk.Owner == pk.Player1 {
		string5 = "| P1 |"
	} else {
		string5 = stringEmpty
	}
	var string1 string
	if pk.Owner != pk.Player1 {
		string1 = "| P2 |"
	} else {
		string1 = stringEmpty
	}
	string3 := fmt.Sprintf("| %d  |", pk.Beads)
	string6 := " ---- "
	return [7]string{string0, string1, stringEmpty, string3, stringEmpty, string5, string6}[line]
}

func CreateDisplayer(player1 models.Player, bowl models.Bowl) Displayer {
	if pb, ok := bowl.(*models.PlayerBowl); ok {
		return &PrintPlayerBowl{pb, player1}
	}
	if pk, ok := bowl.(*models.Kalaha); ok {
		return &PrintKalaha{pk, player1}
	}
	return nil
}
