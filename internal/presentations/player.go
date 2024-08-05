// Package presentations
// Automatic generated
package presentations

import "gitlab.privy.id/privypass/privypass-boilerplate/internal/entity"

type (
	// PlayerQuery parameter
	PlayerQuery struct {
		PlayerName string        `db:"player_name,omitempty" json:"player_name" url:"player_name,omitempty"`
		Position   string        `db:"position,omitempty" json:"position" url:"position,omitempty"`
		Physique   string        `db:"physique,omitempty" json:"physique" url:"physique,omitempty"`
		StatsList  []entity.Stat `db:"stats,omitempty" json:"stats" url:"stats,omitempty"`
		Paging
		PeriodRange
	}

	// PlayerParam input param
	PlayerParam struct {
		Id         string          `db:"id,omitempty" json:"id"`
		PlayerName string          `db:"player_name,omitempty" json:"player_name"`
		Position   string          `db:"position,omitempty" json:"position"`
		Physique   entity.Physique `db:"physique,omitempty" json:"physique"`
		StatsList  []entity.Stat   `db:"stats,omitempty" json:"stats"`
		// AvgStats string	`db:"avgStats,omitempty" json:"avgStats"`
	}

	// PlayerDetail detail response
	PlayerDetail struct {
		// CreatedAt  string `json:"created_at"`
		ID         string          `json:"id"`
		PlayerName string          `json:"player_name"`
		Position   string          `json:"position"`
		Physique   entity.Physique `json:"physique"`
		StatsList  []entity.Stat   `json:"stats"`
	}
)
