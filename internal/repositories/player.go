// Package repositories
// Automatic generated
package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/sync/errgroup"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/entity"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/builderx"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/databasex"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/generator"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/tracer"
)

// Playerer contract of Player
type Playerer interface {
	Store(ctx context.Context, PlayerName string, Position string, Physique entity.Physique, StatsList []entity.Stat) (int64, error)
	Update(ctx context.Context, PlayerName string, Position string, Physique entity.Physique, StatsList []entity.Stat, Id string) (int64, error)
	Delete(ctx context.Context, Id string) (int64, error)
	Counter
	FindOne(ctx context.Context, Id string) (*entity.Player, error)
	Find(ctx context.Context) ([]entity.Player, error)
	FindWithCount(ctx context.Context, param any) ([]entity.Player, int64, error)
}

type player struct {
	db databasex.Adapter
}

// NewPlayer create new instance of Player
func NewPlayer(db databasex.Adapter) Playerer {
	return &player{db: db}
}

// FindOne player
func (r *player) FindOne(ctx context.Context, Id string) (*entity.Player, error) {
	var (
		result entity.Player
		err    error
	)

	ctx = tracer.SpanStart(ctx, "repo.player_find_one")
	defer tracer.SpanFinish(ctx)

	q := `SELECT
			id,
			player_name,
			position
		 FROM players 
		 WHERE id = $1 
		 LIMIT 1`

	err = r.db.QueryRow(ctx, &result, q, Id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	q2 := `SELECT
			height,
			weight,
			age,
			wingspan
		 FROM physique 
		 WHERE id = $1 
		 LIMIT 1`

	err = r.db.QueryRow(ctx, &result.Physique, q2, Id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	q3 := `SELECT
			stats_id,
			points,
			assists,
			rebounds,
			fieldgoalpct,
			threepointpct,
			steals,
			blocks,
			turnovers
		 FROM stats 
		 WHERE id = $1`

	err = r.db.Query(ctx, &result.StatsList, q3, Id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

// Find player
func (r *player) Find(ctx context.Context) ([]entity.Player, error) {
	var (
		result []entity.Player
		err    error
	)

	ctx = tracer.SpanStart(ctx, "repo.player_finds")
	defer tracer.SpanFinish(ctx)

	// wq, err := builderx.StructToMySqlQueryWhere(param, "db")
	// if err != nil {
	// 	tracer.SpanError(ctx, err)
	// 	return nil, err
	// }

	q := `SELECT id, player_name, position FROM players`

	// vals := wq.Values
	// vals = append(vals, wq.Limit, common.PageToOffset(wq.Limit, wq.Page))
	// err = r.db.Query(ctx, &result, fmt.Sprintf(q, wq.Query))
	err = r.db.Query(ctx, &result, q)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	for i := 0; i < len(result); i++ {
		q2 := `SELECT height, weight, age, wingspan FROM physique WHERE id = $1`
		err = r.db.QueryRow(ctx, &result[i].Physique, q2, result[i].Id)
		if err == sql.ErrNoRows {
			return nil, nil
		}

		q3 := `SELECT stats_id, points, assists, rebounds, fieldgoalpct, threepointpct, steals, blocks, turnovers FROM stats WHERE id = $1`
		err = r.db.Query(ctx, &result[i].StatsList, q3, result[i].Id)
		if err == sql.ErrNoRows {
			return nil, nil
		}

	}

	return result, err
}

// Store player
func (r *player) Store(ctx context.Context, PlayerName string, Position string, Physique entity.Physique, StatsList []entity.Stat) (int64, error) {
	var (
		err      error
		affected int64
	)

	generator.Setup(1)
	id := generator.GenerateString()

	ctx = tracer.SpanStart(ctx, "repo.player_store")
	defer tracer.SpanFinish(ctx)

	q := `INSERT INTO players (id, player_name, position) VALUES ($1,$2,$3)`

	_, err = r.db.Exec(ctx, q, id, PlayerName, Position)
	if err != nil {
		return 0, err
	}

	q2 := `INSERT INTO physique (id, height, weight, age, wingspan) VALUES ($1,$2,$3,$4,$5)`

	_, err = r.db.Exec(ctx, q2, id, Physique.Height, Physique.Weight, Physique.Age, Physique.Wingspan)
	if err != nil {
		return 0, err
	}

	for i := 0; i < len(StatsList); i++ {
		q3 := `INSERT INTO stats (
			id,
			stats_id, 
			points, 
			assists, 
			rebounds, 
			fieldGoalPct,
			threePointPct,
			steals,
			blocks,
			turnovers) 
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9, $10)`

		_, err = r.db.Exec(ctx, q3, id, StatsList[i].StatsId, StatsList[i].Points, StatsList[i].Assists, StatsList[i].Rebounds, StatsList[i].FieldGoalPct, StatsList[i].ThreePointPct, StatsList[i].Steals, StatsList[i].Blocks, StatsList[i].Turnovers)
		if err != nil {
			return 0, err
		}
	}

	return affected, nil

}

// Update player data
func (r *player) Update(ctx context.Context, PlayerName string, Position string, Physique entity.Physique, StatsList []entity.Stat, Id string) (int64, error) {
	var (
		err         error
		affected    int64
		statsIdList []string
	)

	ctx = tracer.SpanStart(ctx, "repo.player_update")
	defer tracer.SpanFinish(ctx)

	// query, vals, err := builderx.StructToQueryUpdate(input, where, "players", "db")
	// if err != nil {
	// 	tracer.SpanError(ctx, err)
	// 	return 0, err
	// }

	// // See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	// err = r.db.Transact(ctx, sql.LevelRepeatableRead, func(tx *databasex.DB) error {
	// 	af, err := tx.Exec(ctx, query, vals...)
	// 	affected = af
	// 	return err
	// })

	q := `UPDATE players SET player_name = $1, position = $2 WHERE id = $3`

	_, err = r.db.Exec(ctx, q, PlayerName, Position, Id)
	if err != nil {
		return 0, err
	}

	q2 := `UPDATE physique SET height = $1, weight = $2, age = $3, wingspan = $4 WHERE id = $5`

	affected, err = r.db.Exec(ctx, q2, Physique.Height, Physique.Weight, Physique.Age, Physique.Wingspan, Id)
	if err != nil {
		return 0, err
	}

	q4 := `SELECT stats_id FROM stats`
	err = r.db.Query(ctx, &statsIdList, q4)
	if err != nil {
		return 0, err
	}

	for i := 0; i < len(StatsList); i++ {

		if contains(statsIdList, StatsList[i].StatsId) {
			q3 := `UPDATE stats SET points = $1, assists = $2, rebounds = $3, fieldgoalpct = $4, threepointpct = $5, steals = $6, blocks = $7, turnovers = $8 WHERE stats_id = $9`

			affected, err = r.db.Exec(ctx, q3, StatsList[i].Points, StatsList[i].Assists, StatsList[i].Rebounds, StatsList[i].FieldGoalPct, StatsList[i].ThreePointPct, StatsList[i].Steals, StatsList[i].Blocks, StatsList[i].Turnovers, StatsList[i].StatsId)
			if err != nil {
				return 0, err
			}
		} else {
			q5 := `INSERT INTO stats (id, stats_id, points, assists, rebounds, fieldgoalpct, threepointpct, steals, blocks, turnovers) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

			affected, err = r.db.Exec(ctx, q5, Id, StatsList[i].StatsId, StatsList[i].Points, StatsList[i].Assists, StatsList[i].Rebounds, StatsList[i].FieldGoalPct, StatsList[i].ThreePointPct, StatsList[i].Steals, StatsList[i].Blocks, StatsList[i].Turnovers)
			if err != nil {
				return 0, err
			}
		}

	}

	for i := 0; i < len(statsIdList); i++ {
		valid := false
		for j := 0; j < len(StatsList); j++ {
			if statsIdList[i] == StatsList[j].StatsId {
				valid = true
			}
		}
		if !valid {
			q78 := "DELETE FROM stats WHERE stats_id = $1"

			_, err = r.db.Exec(ctx, q78, statsIdList[i])
			if err != nil {
				return 0, err
			}
		}

	}

	return affected, err
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Delete player from database
func (r *player) Delete(ctx context.Context, Id string) (int64, error) {
	var (
		err error
		// affected         int64
		playersAffected  int64
		physiqueAffected int64
		statsAffected    int64
	)
	ctx = tracer.SpanStart(ctx, "repo.player_delete")
	defer tracer.SpanFinish(ctx)

	// query, vals, err := builderx.StructToQueryDelete(Id, "players", "db", true)
	// if err != nil {
	// 	tracer.SpanError(ctx, err)
	// 	return 0, err
	// }

	// // See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	// err = r.db.Transact(ctx, sql.LevelRepeatableRead, func(tx *databasex.DB) error {
	// 	af, err := tx.Exec(ctx, query, vals...)
	// 	affected = af
	// 	return err
	// })

	q := "DELETE FROM players WHERE id = $1"

	playersAffected, err = r.db.Exec(ctx, q, Id)
	if err != nil {
		return 0, err
	}

	q2 := "DELETE FROM physique WHERE id = $1"

	physiqueAffected, err = r.db.Exec(ctx, q2, Id)
	if err != nil {
		return 0, err
	}

	q3 := "DELETE FROM stats WHERE id = $1"

	statsAffected, err = r.db.Exec(ctx, q3, Id)
	if err != nil {
		return 0, err
	}

	return (playersAffected + physiqueAffected + statsAffected), err
}

// Count player
func (r *player) Count(ctx context.Context, p any) (total int64, err error) {
	ctx = tracer.SpanStart(ctx, "repo.player_count")
	defer tracer.SpanFinish(ctx)

	wq, err := builderx.StructToMySqlQueryWhere(p, "db")
	if err != nil {
		tracer.SpanError(ctx, err)
		return
	}

	q := fmt.Sprintf(`
		SELECT
        	COUNT(id) AS jumlah
		FROM players %s `, wq.Query)

	err = r.db.QueryRow(ctx, &total, q, wq.Values...)
	if err != nil {
		tracer.SpanError(ctx, err)
		err = err
		return
	}

	return
}

// FindWithCount find player with count
func (r *player) FindWithCount(ctx context.Context, param any) ([]entity.Player, int64, error) {

	var (
		cl    []entity.Player
		count int64
	)

	ctx = tracer.SpanStart(ctx, "repo.player_with_count")
	defer tracer.SpanFinish(ctx)

	group, newCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		l, err := r.Find(newCtx)
		cl = l
		return err
	})
	group.Go(func() error {
		c, err := r.Count(ctx, param)
		count = c
		return err
	})

	err := group.Wait()

	return cl, count, err
}
