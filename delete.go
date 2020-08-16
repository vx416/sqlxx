package sqlxx

type DeleteCtx struct {
}

func (ctx DeleteCtx) SetupContext(exec *Executor) (queryStr string, args []interface{}, err error) {
	query := exec.QueryOpts

	usingTable := make([]string, 0)
	if exec.data != nil {
		usingTable = exec.data.([]string)
	}
	queryStr, args, err = exec.db.GetBuilder().BuildDelete(exec.table, usingTable, query)
	if err != nil {
		return "", []interface{}{}, err
	}
	return queryStr, args, nil
}

func (ctx DeleteCtx) ExecHandle(exec *Executor, query string, args ...interface{}) error {
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
