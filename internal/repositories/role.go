// Package repositories
// Automatic generated
package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/sync/errgroup"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/common"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/entity"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/builderx"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/databasex"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/tracer"
)

// Roler contract of Role
type Roler interface {
	Update(ctx context.Context, Name string, Code string) (int64, error)
	Deleter
	Counter
	FindOne(ctx context.Context, Code string) (*entity.Role, error)
	Find(ctx context.Context, param any) ([]entity.Role, error)
	FindWithCount(ctx context.Context, param any) ([]entity.Role, int64, error)
	Store(ctx context.Context, Email string, Password string, Code string) error
}

type role struct {
	db databasex.Adapter
}

// NewRole create new instance of Role
func NewRole(db databasex.Adapter) Roler {
	return &role{db: db}
}

// FindOne role
func (r *role) FindOne(ctx context.Context, Code string) (*entity.Role, error) {
	var (
		result entity.Role
		err    error
	)

	ctx = tracer.SpanStart(ctx, "repo.role_find_one")
	defer tracer.SpanFinish(ctx)

	// wq, err := builderx.StructToMySqlQueryWhere(param, "db")
	// if err != nil {
	// 	tracer.SpanError(ctx, err)
	// 	return nil, err
	// }

	q := `SELECT 
			code,
			name,
			created_at,
			updated_at,
			deleted_at
		 FROM roles WHERE code = $1 LIMIT 1`

	err = r.db.QueryRow(ctx, &result, q, Code)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

// Find role
func (r *role) Find(ctx context.Context, param any) ([]entity.Role, error) {
	var (
		result []entity.Role
		err    error
	)

	ctx = tracer.SpanStart(ctx, "repo.role_finds")
	defer tracer.SpanFinish(ctx)

	wq, err := builderx.StructToMySqlQueryWhere(param, "db")
	if err != nil {
		tracer.SpanError(ctx, err)
		return nil, err
	}

	q := `SELECT 
			code,
			name,
			created_at,
			updated_at,
			deleted_at
		 FROM roles %s LIMIT ? OFFSET ? `

	vals := wq.Values
	vals = append(vals, wq.Limit, common.PageToOffset(wq.Limit, wq.Page))
	err = r.db.Query(ctx, &result, fmt.Sprintf(q, wq.Query), vals...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return result, err
}

// Store role
func (r *role) Store(ctx context.Context, Email string, Password string, Code string) error {
	var (
		err error
	)

	ctx = tracer.SpanStart(ctx, "repo.role_store")
	defer tracer.SpanFinish(ctx)

	// np := &param
	// param = *np
	// query, vals, err := builderx.StructToQueryInsert(param, "roles", "db")
	// if err != nil {
	// 	tracer.SpanError(ctx, err)
	// 	return err
	// }

	// // See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	// err = r.db.Transact(ctx, sql.LevelRepeatableRead, func(tx *databasex.DB) error {
	// 	af, err := tx.Exec(ctx, query, vals...)
	// 	affected = af
	// 	return err
	// })

	q := `INSERT INTO roles (email, password, code) VALUES ($1,$2, $3)`

	_, err = r.db.Exec(ctx, q, Email, Password, Code)
	if err != nil {
		return err
	}

	return nil
}

// Update role data
func (r *role) Update(ctx context.Context, Email string, Code string) (int64, error) {
	var (
		err      error
		affected int64
	)

	ctx = tracer.SpanStart(ctx, "repo.role_update")
	defer tracer.SpanFinish(ctx)

	// query, vals, err := builderx.StructToQueryUpdate(input, where, "roles", "db")
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

	q := `UPDATE roles SET code = $1 WHERE email = $2`

	affected, err = r.db.Exec(ctx, q, Code, Email)
	if err != nil {
		return 0, err
	}

	return affected, nil

}

// Delete role from database
func (r *role) Delete(ctx context.Context, param any) (int64, error) {
	var (
		err      error
		affected int64
	)
	ctx = tracer.SpanStart(ctx, "repo.role_delete")
	defer tracer.SpanFinish(ctx)

	query, vals, err := builderx.StructToQueryDelete(param, "roles", "db", true)
	if err != nil {
		tracer.SpanError(ctx, err)
		return 0, err
	}

	// See https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels.
	err = r.db.Transact(ctx, sql.LevelRepeatableRead, func(tx *databasex.DB) error {
		af, err := tx.Exec(ctx, query, vals...)
		affected = af
		return err
	})

	return affected, err
}

// Count role
func (r *role) Count(ctx context.Context, p any) (total int64, err error) {
	ctx = tracer.SpanStart(ctx, "repo.role_count")
	defer tracer.SpanFinish(ctx)

	wq, err := builderx.StructToMySqlQueryWhere(p, "db")
	if err != nil {
		tracer.SpanError(ctx, err)
		return
	}

	q := fmt.Sprintf(`
		SELECT
        	COUNT(id) AS jumlah
		FROM roles %s `, wq.Query)

	err = r.db.QueryRow(ctx, &total, q, wq.Values...)
	if err != nil {
		tracer.SpanError(ctx, err)
		err = err
		return
	}

	return
}

// FindWithCount find role with count
func (r *role) FindWithCount(ctx context.Context, param any) ([]entity.Role, int64, error) {

	var (
		cl    []entity.Role
		count int64
	)

	ctx = tracer.SpanStart(ctx, "repo.role_with_count")
	defer tracer.SpanFinish(ctx)

	group, newCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		l, err := r.Find(newCtx, param)
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
