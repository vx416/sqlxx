package sqlxx

type Execr interface {
	Do() error
	LastInsertID() int64
	RowsAffected() int64
}

type ExecContext interface {
	SetupContext(exec *Executor) (query string, args []interface{}, err error)
	ExecHandle(exec *Executor, query string, args ...interface{}) error
}

func NewExec(db *DB, execCtx ExecContext, table string, data interface{}) *Executor {
	exec := &Executor{
		execContext: execCtx,
		table:       table,
		data:        data,
		db:          db,
	}
	return exec
}

type queryHandler func(exec *Executor)
type Executor struct {
	table        string
	data         interface{}
	execContext  ExecContext
	db           *DB
	QueryOpts    *QueryOptions
	lastInsertID int64
	rowsAffected int64
}

func (exec *Executor) Do() error {
	query, args, err := exec.execContext.SetupContext(exec)
	if err != nil {
		return err
	}
	query = exec.db.Rebind(query)

	if err := exec.execContext.ExecHandle(exec, query, args...); err != nil {
		return err
	}
	return nil
}

func (exec *Executor) LastInsertID() int64 {
	return exec.lastInsertID
}

func (exec *Executor) RowsAffected() int64 {
	return exec.rowsAffected
}
