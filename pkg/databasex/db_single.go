package databasex

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/tracer"

	"github.com/jmoiron/sqlx"
)

var (
	// check in runtime implement Databaser
	_ Adapter = (*DB)(nil)
)

type DB struct {
	db *sqlx.DB
	//instanceID string
	tx   *sqlx.Tx
	conn *sqlx.Conn // the Conn of the Tx, when tx != nil
	//opts       sql.TxOptions // valid when tx != nil
	reaMode bool
	dbName  string
}

type NamedStmt struct {
	dbName string `json:"db_name"`
	db     *sqlx.NamedStmt
}

func New(db *sqlx.DB, readMode bool, sbName string) *DB {
	return &DB{
		db:      db,
		reaMode: readMode,
		dbName:  sbName,
	}
}

func (db *DB) Ping() error {
	return db.db.Ping()
}

func (db *DB) InTransaction() bool {
	return db.tx != nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.db.Close()
}

// Exec executes a SQL statement and returns the number of rows it affected.
func (db *DB) Exec(ctx context.Context, query string, args ...any) (_ int64, err error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "exec",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
	)
	defer tracer.SpanFinish(ctx)
	if db.reaMode {
		return 0, fmt.Errorf("database mode read only")
	}

	res, err := db.execResult(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("rows affected: %v", err)
	}

	return n, nil
}

// execResult executes a SQL statement and returns a sql.Result.
func (db *DB) execResult(ctx context.Context, query string, args ...any) (res sql.Result, err error) {
	if db.tx != nil {
		return db.tx.ExecContext(ctx, query, args...)
	}

	return db.db.ExecContext(ctx, query, args...)
}

// Query runs the DB query.
func (db *DB) Query(ctx context.Context, dst any, query string, args ...any) error {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "query",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
	)
	defer tracer.SpanFinish(ctx)
	if db.tx != nil {
		return db.tx.SelectContext(ctx, dst, query, args...)
	}

	return db.db.SelectContext(ctx, dst, query, args...)
}

// QueryRow runs the query and returns a single row.
func (db *DB) QueryRow(ctx context.Context, dst any, query string, args ...any) error {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "query_row",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
	)
	defer tracer.SpanFinish(ctx)

	if db.tx != nil {
		return db.tx.GetContext(ctx, dst, query, args...)
	}

	return db.db.GetContext(ctx, dst, query, args...)
}

// QueryX runs the DB query.
func (db *DB) QueryX(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "queryx",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
	)
	defer tracer.SpanFinish(ctx)
	if db.tx != nil {
		return db.tx.QueryContext(ctx, query, args...)
	}

	return db.db.QueryContext(ctx, query, args...)
}

// QueryRowX runs the query and returns a single row.
func (db *DB) QueryRowX(ctx context.Context, query string, args ...any) *sql.Row {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "query_rowx",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
	)
	defer tracer.SpanFinish(ctx)
	if db.tx != nil {
		return db.tx.QueryRowContext(ctx, query, args...)
	}

	return db.db.QueryRowContext(ctx, query, args...)
}

func (db *DB) PrepareNamedContext(ctx context.Context, query string) (*NamedStmt, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "prepare_named_context",
		tracer.WithResourceNameOptions(query),
		tracer.WithOptions("sql.query", query),
	)

	defer tracer.SpanFinish(ctx)
	if db.InTransaction() {
		namedStmt, err := db.tx.PrepareNamedContext(ctx, query)
		if err != nil {
			return nil, err
		}

		return &NamedStmt{
			dbName: db.dbName,
			db:     namedStmt,
		}, nil
	}

	namedStmt, err := db.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}

	return &NamedStmt{
		dbName: db.dbName,
		db:     namedStmt,
	}, nil
}

// Transact executes the given function in the context of a SQL transaction at
// the given isolation level
func (db *DB) Transact(ctx context.Context, iso sql.IsolationLevel, txFunc func(*DB) error) (err error) {
	ctx = tracer.DBSpanStartWithOption(ctx, db.dbName, "transaction")
	defer tracer.SpanFinish(ctx)
	if db.reaMode {
		return fmt.Errorf("database mode read only")
	}

	// For the levels which require retry, see
	// https://www.postgresql.org/docs/11/transaction-iso.html.
	opts := &sql.TxOptions{Isolation: iso}

	return db.transact(ctx, opts, txFunc)
}

func (db *DB) transact(ctx context.Context, opts *sql.TxOptions, txFunc func(*DB) error) (err error) {
	if db.InTransaction() {
		return errors.New("db transact function was called on a DB already in a transaction")
	}

	conn, err := db.db.Connx(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	tx, err := conn.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("tx begin: %w", err)
	}

	//defer func() {
	//	if p := recover(); p != nil {
	//		tx.Rollback()
	//	} else if err != nil {
	//		tx.Rollback()
	//	} else {
	//		if txErr := tx.Commit(); txErr != nil {
	//			err = fmt.Errorf("tx commit: %w", txErr)
	//		}
	//	}
	//}()

	dbtx := New(db.db, false, db.dbName)
	dbtx.tx = tx
	dbtx.conn = conn
	//dbtx.opts = *opts

	if err := txFunc(dbtx); err != nil {
		tx.Rollback()
		return fmt.Errorf("fn(tx): %w", err)
	}

	return tx.Commit()
}

func (db *DB) DBRead() Adapter {
	return db
}

func (db *DB) DBWrite() Adapter {
	return db
}

func (db *DB) Rebind(ctx context.Context, q string) string {
	return db.db.Rebind(q)
}

func (db *DB) BeginTx(ctx context.Context, iso sql.IsolationLevel) Adapter {
	db.tx = db.db.MustBeginTx(ctx, &sql.TxOptions{Isolation: iso})
	return db
}

func (db *DB) Commit(ctx context.Context) error {
	if !db.InTransaction() {
		return errors.New("db not in transaction mode")
	}
	defer func() {
		db.tx = nil
	}()
	err := db.tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Rollback(ctx context.Context) error {
	if !db.InTransaction() {
		return errors.New("db not in transaction mode")
	}
	defer func() {
		db.tx = nil
	}()
	err := db.tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func (n *NamedStmt) QueryxContext(ctx context.Context, arg interface{}) (*sqlx.Rows, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, n.dbName, "prepare_named_context.queryx_context",
		tracer.WithResourceNameOptions(n.db.QueryString),
		tracer.WithOptions("sql.query", n.db.QueryString),
	)

	defer tracer.SpanFinish(ctx)

	return n.db.QueryxContext(ctx, arg)
}

func (n *NamedStmt) QueryRowxContext(ctx context.Context, arg interface{}) *sqlx.Row {
	ctx = tracer.DBSpanStartWithOption(ctx, n.dbName, "prepare_named_context.queryx_context",
		tracer.WithResourceNameOptions(n.db.QueryString),
		tracer.WithOptions("sql.query", n.db.QueryString),
	)

	defer tracer.SpanFinish(ctx)
	return n.db.QueryRowxContext(ctx, arg)
}

func (n *NamedStmt) ExecContext(ctx context.Context, arg interface{}) (sql.Result, error) {
	ctx = tracer.DBSpanStartWithOption(ctx, n.dbName, "prepare_named_context.exec_context",
		tracer.WithResourceNameOptions(n.db.QueryString),
		tracer.WithOptions("sql.query", n.db.QueryString),
	)
	defer tracer.SpanFinish(ctx)

	return n.db.ExecContext(ctx, arg)
}

func (n NamedStmt) GetContext(ctx context.Context, dest interface{}, arg interface{}) error {
	ctx = tracer.DBSpanStartWithOption(ctx, n.dbName, "prepare_named_context.get_context",
		tracer.WithResourceNameOptions(n.db.QueryString),
		tracer.WithOptions("sql.query", n.db.QueryString),
	)
	defer tracer.SpanFinish(ctx)

	return n.db.GetContext(ctx, dest, arg)
}
