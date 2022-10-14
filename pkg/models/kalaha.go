package models

type Kalaha struct {
	Name string
	BaseBowl
}

func (k *Kalaha) PassBeads(player Player, beads uint) Player {
	if k.Owner == player && beads > 0 {
		k.Beads++
		if beads == 1 {
			return player
		}
		return (k.TheNext).PassBeads(player, beads-1)
	}
	return (k.TheNext).PassBeads(player, beads)

}
