package sqlxx

import (
	"context"
	"errors"
	"sync"

	"github.com/jmoiron/sqlx"
)

type (
	MasterSlaveKey    struct{}
	IsolationLevelKey struct{}
)

func WithSlave(ctx context.Context) context.Context {
	return context.WithValue(ctx, MasterSlaveKey{}, 1)
}

func WithMaster(ctx context.Context) context.Context {
	return context.WithValue(ctx, MasterSlaveKey{}, 2)
}

func IsReadOnly(ctx context.Context) bool {
	val, ok := ctx.Value(MasterSlaveKey{}).(int)
	return ok && val == 1
}

func IsMaster(ctx context.Context) bool {
	val, ok := ctx.Value(MasterSlaveKey{}).(int)
	return ok && val == 2
}

type Policy interface {
	Get(ctx context.Context) (*sqlx.DB, error)
}

func NewRRCluster(masters []*sqlx.DB, slaves []*sqlx.DB) *Cluster {
	return &Cluster{
		masters: &RoundRubinPolicy{dbs: masters},
		slaves:  &RoundRubinPolicy{dbs: masters},
	}
}

type Cluster struct {
	masters Policy
	slaves  Policy
}

func (c Cluster) GetDB(ctx context.Context) (*sqlx.DB, error) {
	if IsReadOnly(ctx) {
		return c.slaves.Get(ctx)
	}
	return c.masters.Get(ctx)
}

type RoundRubinPolicy struct {
	dbs       []*sqlx.DB
	lock      sync.Mutex
	currIndex int
}

func (po *RoundRubinPolicy) Get(ctx context.Context) (*sqlx.DB, error) {
	po.lock.Lock()
	defer po.lock.Unlock()
	if len(po.dbs) == 0 {
		return nil, errors.New("db pool is not set up")
	}

	if po.currIndex >= len(po.dbs) {
		po.currIndex = 0
	}

	db := po.dbs[po.currIndex]
	po.currIndex++
	return db, nil
}
