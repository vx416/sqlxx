package sqlxx

import (
	"github.com/vicxu416/sqlxx/sqlbuilder/query"
)

func NewQueryOpts() *QueryOptions {
	return &QueryOptions{
		Querier: query.NewQuery(),
	}
}

type QueryOptions struct {
	query.Querier
}
