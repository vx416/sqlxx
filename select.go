package sqlxx

type SelectCtx struct {
	get bool
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

	if ctx.get == true {
		return exec.db.proxyDB.Get(exec.data, query, args...)
	}
	return exec.db.proxyDB.Select(exec.data, query, args...)
}

// type GetCtx struct {
// }

// func (ctx FindCtx) SetupContext(exec *Executor) (queryStr string, args []interface{}, err error) {
// 	query := exec.QueryOpts
// 	queryStr, args, err = exec.db.GetBuilder().BuildSelect(exec.table, query)
// 	if err != nil {
// 		return "", []interface{}{}, err
// 	}

// 	return queryStr, args, nil
// }
// func (ctx FindCtx) ExecHandle(exec *Executor, query string, args ...interface{}) error {
// 	return exec.db.proxyDB.Get(exec.data, query, args...)
// }
