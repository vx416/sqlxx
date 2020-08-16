package sqlxx

type SelectCtx struct {
}

func (ctx SelectCtx) SetupContext(exec *Executor) (queryStr string, args []interface{}, err error) {
	query := exec.QueryOpts
	queryStr, args, err = exec.db.GetBuilder().BuildSelect(exec.table, query)
	if err != nil {
		return "", []interface{}{}, err
	}

	return queryStr, args, nil
}
func (ctx SelectCtx) ExecHandle(exec *Executor, query string, args ...interface{}) error {
	return exec.db.proxyDB.Select(exec.data, query, args...)
}
