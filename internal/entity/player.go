// Package entity
// Automatic generated
package entity

// Player entity
type Player struct {
	Id string `db:"id,omitempty" json:"id"`
	// CreatedAt  time.Time `db:"created_at,omitempty" json:"created_at"`
	PlayerName string   `db:"player_name,omitempty" json:"player_name"`
	Position   string   `db:"position,omitempty" json:"position"`
	Physique   Physique `db:"physique,omitempty" json:"physique"`
	StatsList  []Stat   `db:"stats,omitempty" json:"stats"`
}

// type StatsList struct {
// 	Stats []Stat
// }

type Physique struct {
	Height   string `db:"height,omitempty" json:"height"`
	Weight   string `db:"weight,omitempty" json:"weight"`
	Age      string `db:"age,omitempty" json:"age"`
	Wingspan string `db:"wingspan,omitempty" json:"wingspan"`
}

type Stat struct {
	StatsId       string `db:"stats_id,omitempty" json:"stats_id"`
	Points        string `db:"points,omitempty" json:"points"`
	Assists       string `db:"assists,omitempty" json:"assists"`
	Rebounds      string `db:"rebounds,omitempty" json:"rebounds"`
	FieldGoalPct  string `db:"fieldgoalpct,omitempty" json:"fieldgoalpct"`
	ThreePointPct string `db:"threepointpct,omitempty" json:"threepointpct"`
	Steals        string `db:"steals,omitempty" json:"steals"`
	Blocks        string `db:"blocks,omitempty" json:"blocks"`
	Turnovers     string `db:"turnovers,omitempty" json:"turnovers"`
}
