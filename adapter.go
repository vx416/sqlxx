package sqlxx

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	TxKey struct{}
)

func NewWith(sqlxDB *sqlx.DB) *Sqlxx {
	return &Sqlxx{
		db: &DB{
			Cluster: NewRRCluster([]*sqlx.DB{sqlxDB}, []*sqlx.DB{sqlxDB}),
		},
	}
}

func NewWithCluster(masters, slaves []*sqlx.DB) *Sqlxx {
	return &Sqlxx{
		db: &DB{
			Cluster: NewRRCluster(masters, slaves),
		},
	}
}

type Sqlxx struct {
	db *DB
}

func (adapter *Sqlxx) GetDB(ctx context.Context) *DB {
	txDB := adapter.getTx(ctx)
	if txDB != nil {
		return txDB
	}
	return adapter.db
}

func (adapter *Sqlxx) ExecuteTx(ctx context.Context, fn func(txCtx context.Context) error, txOpt *sql.TxOptions) error {
	txDB, err := adapter.db.Begin(ctx, txOpt)
	if err != nil {
		return err
	}

	txCtx := adapter.withTx(ctx, txDB)
	defer func() {
		if pErr := recover(); pErr != nil {
			txDB.Rollback(ctx)
		}
	}()
	var callbackErr, txErr error
	callbackErr = fn(ctx)
	if callbackErr != nil {
		txErr = txDB.Rollback(txCtx)
	} else {
		if txErr = txDB.Commit(txCtx); txErr != nil {
			txErr = txDB.Rollback(txCtx)
		}
	}

	if txErr != nil {
		if callbackErr != nil {
			txErr = errors.Wrapf(txErr, "callback error:%+v", callbackErr)
		}
		return txErr
	}
	return callbackErr
}

func (adapter *Sqlxx) ViewTx(ctx context.Context, fn func(ctx context.Context) error, txOpt *sql.TxOptions) error {
	ctx = WithSlave(ctx)
	return adapter.ExecuteTx(ctx, fn, txOpt)
}

func (adapter *Sqlxx) HasTx(ctx context.Context) bool {
	txDB := adapter.getTx(ctx)
	return txDB != nil
}

func (adapter *Sqlxx) withTx(ctx context.Context, db *DB) context.Context {
	return context.WithValue(ctx, TxKey{}, db)
}

func (adapter *Sqlxx) getTx(ctx context.Context) *DB {
	txDB, ok := ctx.Value(TxKey{}).(*DB)
	if ok && txDB != nil {
		return txDB
	}

	return nil
}
