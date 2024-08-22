package internal

import (
	"github.com/sirupsen/logrus"
	"mancala/pkg"
	"mancala/pkg/models"
	"mancala/pkg/utils"
)

func Train(n uint) *models.AIPlayer {
	logrus.Info("Training Starting")
	player1 := &models.AIPlayer{Alpha: 0.5, Epsilon: 0.1, BasePlayer: &models.BasePlayer{}}
	player2 := &models.AIPlayer{Alpha: 0.5, Epsilon: 0.1, BasePlayer: &models.BasePlayer{Opponent: player1}}
	player1.Opponent = player2

	for i := uint(1); i < n; i++ {
		ge := pkg.Initialize(player1, player2)
		last := map[models.Player]utils.Pair{}
		for ge.PlayerTurn.CanPlay() {
			state := ge.GetState()
			action := ge.PlayerTurn.(*models.AIPlayer).ChooseAction(state)
			last[ge.PlayerTurn] = utils.PairFrom(state, action)
			err := ge.Play(action)
			if err != nil {
				panic("error in move: " + err.Error())
			}

			if ge.PlayerTurn.CanPlay() {
				if pair, ok := last[ge.PlayerTurn]; ok {
					ge.Player1.(*models.AIPlayer).Update(pair, ge.GetState(), 0)
					ge.Player2.(*models.AIPlayer).Update(pair, ge.GetState(), 0)
				}
			} else {
				// game is finished
				points1 := ge.GetPoints(ge.Player1)
				points2 := ge.GetPoints(ge.Player2)
				var winner models.Player
				var loser models.Player
				if points1 > points2 {
					winner = ge.Player1
					loser = ge.Player2
				} else {
					winner = ge.Player2
					loser = ge.Player1
				}
				ge.Player1.(*models.AIPlayer).Update(last[winner], ge.GetState(), 1)
				ge.Player2.(*models.AIPlayer).Update(last[winner], ge.GetState(), 1)

				ge.Player1.(*models.AIPlayer).Update(last[loser], ge.GetState(), -1)
				ge.Player2.(*models.AIPlayer).Update(last[loser], ge.GetState(), -1)

			}
		}
		if i%100 == 0 {
			logrus.WithField("i", i).Info("Training Game Finished")
		}
	}
	logrus.Info("Training Finished")
	return player2

}
