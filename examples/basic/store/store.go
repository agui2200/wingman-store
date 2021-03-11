// Code generated by wingman-store, DO NOT EDIT.
package store

import (
	"bytes"
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/agui2200/wingman-store/examples/basic/store/migrate"
	_ "github.com/agui2200/wingman-store/examples/basic/store/runtime"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
)

type databaseKey int

var migrationDataHook []func(ctx context.Context) error

var (
	ErrUnsupportedDrive = errors.New("unsupported driver")
	ErrDatabaseNotFound = errors.New("database connection not found")
)

const (
	ctxDatabaseKey databaseKey = iota
	ctxDatabaseTxKey
)

func New(ctx context.Context, driverName, dataSourceName string, debug bool) (context.Context, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		var options []Option
		if debug {
			dbgDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
				log.Println("exec sql:", i)
			})
			options = append(options, Driver(dbgDrv))
		}
		c := NewClient(options...)
		return context.WithValue(ctx, ctxDatabaseKey, c), nil
	default:
		return nil, errors.WithMessage(ErrUnsupportedDrive, "driveName:"+driverName)
	}
}

func NewTx(ctx context.Context) (context.Context, error) {
	db := WithContext(ctx)
	if db != nil {
		c, err := db.Tx(ctx)
		if err != nil {
			return nil, err
		}
		return context.WithValue(ctx, ctxDatabaseTxKey, c), nil
	}
	return nil, ErrUnsupportedDrive
}

func WithContext(ctx context.Context) *Client {
	// 要是有事务执行，优先用事务
	if ctx.Value(ctxDatabaseTxKey) != nil {
		if c, ok := ctx.Value(ctxDatabaseTxKey).(*Tx); ok {
			return c.Client()
		}
	}
	if ctx.Value(ctxDatabaseKey) != nil {
		if c, ok := ctx.Value(ctxDatabaseKey).(*Client); ok {
			return c
		}
	}

	return nil
}

func WithTxContext(ctx context.Context) *Tx {
	if ctx.Value(ctxDatabaseTxKey) != nil {
		if c, ok := ctx.Value(ctxDatabaseTxKey).(*Tx); ok {
			return c
		}
	}
	return nil
}

func CloseDatabase(ctx context.Context) error {
	if c, ok := ctx.Value(ctxDatabaseKey).(*Client); ok {
		return c.Close()
	}
	return nil
}

func Migration(ctx context.Context, debug bool) error {
	buffer := bytes.NewBuffer([]byte{})
	err := WithContext(ctx).Schema.WriteTo(ctx, buffer)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadAll(buffer)
	if err != nil {
		return err
	}
	if debug {
		log.Printf("migrationDatabase: \r\n%s", buf)
	}
	err = WithContext(ctx).Schema.Create(ctx, migrate.WithGlobalUniqueID(true))
	if err != nil {
		return err
	}

	err = WithTx(ctx, func(ctx context.Context) error {
		for _, f := range migrationDataHook {
			err := f(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx, err := NewTx(ctx)
	if err != nil {
		return err
	}
	tx := WithTxContext(ctx)
	if tx == nil {
		return ErrDatabaseNotFound
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			// handler panic
			err = fmt.Errorf("transaction panic: %v", v)

		}
	}()
	if err := fn(ctx); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = errors.Wrapf(err, "rolling back transaction: %v", rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return errors.Wrapf(err, "committing transaction: %v", err)
	}
	return nil
}
