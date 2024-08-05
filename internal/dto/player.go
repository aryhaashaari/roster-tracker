package dto

import (
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/entity"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/presentations"
)

func PlayerToResponse(src entity.Player) presentations.PlayerDetail {
	x := presentations.PlayerDetail{
		ID:         src.Id,
		PlayerName: src.PlayerName,
		Position:   src.Position,
		Physique:   src.Physique,
		StatsList:  src.StatsList,
	}

	// if !src.CreatedAt.IsZero() {
	// 	x.CreatedAt = src.CreatedAt.Format(consts.LayoutDateTimeFormat)
	// }
	//
	//if !src.UpdatedAt.IsZero() {
	//	x.UpdatedAt = src.UpdatedAt.Format(consts.LayoutDateTimeFormat)
	//}
	//
	//if !src.DeletedAt.IsZero() {
	//	x.DeletedAt = src.DeletedAt.Format(consts.LayoutDateTimeFormat)
	//}

	return x
}

func PlayersToResponse(inputs []entity.Player) []presentations.PlayerDetail {
	var (
		result = []presentations.PlayerDetail{}
	)

	for i := 0; i < len(inputs); i++ {
		result = append(result, PlayerToResponse(inputs[i]))
	}

	return result
}

func DetailedPlayerToResponse(src *entity.Player) presentations.PlayerDetail {
	x := presentations.PlayerDetail{
		ID:         src.Id,
		PlayerName: src.PlayerName,
		Position:   src.Position,
		Physique:   src.Physique,
		StatsList:  src.StatsList,
	}

	return x
}
