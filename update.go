package sqlxx

type UpdateCtx struct {
}

func (ctx UpdateCtx) SetupContext(exec *Executor) (queryStr string, args []interface{}, err error) {
	query := exec.QueryOpts
	whereStr, args, err := exec.db.GetBuilder().BuildUpdate(exec.table, exec.data, query)
	if err != nil {
		return "", []interface{}{}, err
	}
	return whereStr, args, nil
}

func (ctx UpdateCtx) ExecHandle(exec *Executor, query string, args ...interface{}) error {
	res, err := exec.db.proxyDB.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	exec.rowsAffected = rows
	return nil
}

type UpdateStructCtx struct {
}

func (ctx UpdateStructCtx) SetupContext(exec *Executor) (queryStr string, args []interface{}, err error) {
	query, args, err := exec.db.GetBuilder().BuildUpdateStruct(exec.table, exec.data)
	if err != nil {
		return "", []interface{}{}, err
	}
	return query, args, nil
}

func (ctx UpdateStructCtx) ExecHandle(exec *Executor, query string, args ...interface{}) error {
	res, err := exec.db.proxyDB.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	exec.rowsAffected = rows
	return nil
}
