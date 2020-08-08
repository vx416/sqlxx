package sqlxx

type QueryContext interface {
	SetupContext(exec *Executor) error
	QueryHandle(exec *Executor, query string, args ...interface{}) error
}

func NewExec(queryCtx QueryContext, db *DB, table string, data interface{}) *Executor {
	return &Executor{
		queryContext: queryCtx,
		db:           db,
		table:        table,
		data:         data,
	}
}

type Executor struct {
	table        string
	data         interface{}
	queryContext QueryContext
	query        string
	args         []interface{}
	db           *DB
	LastInsertID int64
	RowsAffected int64
}

type queryHandler func(exec *Executor)

func (exec *Executor) addQueryArgs(query string, arg ...interface{}) {
	exec.args = append(exec.args, arg...)
	exec.query = query
}

func (exec *Executor) Do() error {
	if err := exec.queryContext.SetupContext(exec); err != nil {
		return err
	}
	query := exec.db.Rebind(exec.query)

	if err := exec.queryContext.QueryHandle(exec, query, exec.args...); err != nil {
		return err
	}
	return nil
}
