package query

type Querier interface {
	Select(query string) Querier
	Group(query string) Querier
	Having(query string) Querier
	Lock(query string) Querier
	Offset(arg int) Querier
	Limit(arg int) Querier
	Join(query string) Querier
	OrderBy(column string, mode string) Querier

	Where(query string, args ...interface{}) Querier
	In(column string, args interface{}) Querier
	Or(query string, args ...interface{}) Querier
	AndStruct(args interface{}) Querier
	OrStruct(args interface{}) Querier

	WhereQueryArgs() (string, []interface{}, error)
	SelectQueriesArgs() (map[string]string, []interface{}, error)
}

func NewQuery() Querier {
	return &Query{
		Selector: NewSelect(),
		Wherer:   NewWhere(),
	}
}

type Query struct {
	Selector *Select
	Wherer   *Where
}

func (q *Query) Select(query string) Querier {
	q.Selector.Select(query)
	return q
}

func (q *Query) Group(query string) Querier {
	q.Selector.Group(query)
	return q
}

func (q *Query) Having(query string) Querier {
	q.Selector.Having(query)
	return q
}

func (q *Query) Lock(query string) Querier {
	q.Selector.Lock(query)
	return q
}

func (q *Query) Offset(arg int) Querier {
	q.Selector.Offset(arg)
	return q
}

func (q *Query) Limit(arg int) Querier {
	q.Selector.Limit(arg)
	return q
}

func (q *Query) Join(query string) Querier {
	q.Selector.Join(query)
	return q
}

func (q *Query) OrderBy(column string, mode string) Querier {
	q.Selector.OrderBy(column, mode)
	return q
}

func (q *Query) Where(query string, args ...interface{}) Querier {
	q.Wherer.Where(query, args...)
	return q
}

func (q *Query) In(query string, args interface{}) Querier {
	q.Wherer.In(query, args)
	return q
}
func (q *Query) Or(query string, args ...interface{}) Querier {
	q.Wherer.Or(query, args...)
	return q
}
func (q *Query) AndStruct(arg interface{}) Querier {
	q.Wherer.AndStruct(arg)
	return q

}
func (q *Query) OrStruct(arg interface{}) Querier {
	q.Wherer.OrStruct(arg)
	return q
}

func (q *Query) WhereQueryArgs() (string, []interface{}, error) {
	return q.Wherer.QueryArgs()
}

func (q *Query) SelectQueriesArgs() (map[string]string, []interface{}, error) {
	return q.Selector.QueriesArgs()
}
