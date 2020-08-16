package query

import (
	"strconv"
	"strings"
)

type SelectPart string

const (
	SELECT SelectPart = "SELECT"
	JOIN   SelectPart = "JOIN"
	WHERE  SelectPart = "WHERE"
	GROUP  SelectPart = "GROUP BY"
	HAVING SelectPart = "HAVING"
	ORDER  SelectPart = "ORDER BY"
	LIMIT  SelectPart = "LIMIT"
	OFFSET SelectPart = "OFFSET"
	FOR    SelectPart = "FOR"
)

var partsSort = []SelectPart{GROUP, HAVING, LIMIT, OFFSET, FOR}

type rawSQL struct {
	query string
	args  []interface{}
}

func NewSelect() *Select {
	return &Select{
		parts: make(map[SelectPart]string),
	}
}

type Select struct {
	parts   map[SelectPart]string
	lastErr error
}

func (sel *Select) Select(query string) *Select {
	sel.concateStr(SELECT, query, ", ")
	return sel
}

func (sel *Select) Join(query string) *Select {
	sel.parts[JOIN] = query
	return sel
}

func (sel *Select) Group(query string) *Select {
	sel.concateStr(GROUP, query, ", ")
	return sel
}

func (sel *Select) Having(query string) *Select {
	sel.concateStr(HAVING, query, ", ")
	return sel
}

func (sel *Select) OrderBy(column string, mode string) *Select {
	sel.concateStr(ORDER, column+" "+mode, ", ")
	return sel
}

func (sel *Select) Limit(query int) *Select {
	sel.parts[LIMIT] = strconv.Itoa(query)
	return sel
}

func (sel *Select) Offset(query int) *Select {
	sel.parts[OFFSET] = strconv.Itoa(query)
	return sel
}

func (sel *Select) Lock(query string) *Select {
	sel.parts[FOR] = query
	return sel
}

func (sel *Select) LastErr() error {
	return nil
}

func (sel *Select) QueriesArgs() (map[string]string, []interface{}, error) {
	result := make(map[string]string)

	result["select"] = sel.parts[SELECT]
	result["options"] = sel.buildPartsStr()
	result["join"] = sel.parts[JOIN]
	return result, []interface{}{}, nil
}

func (sel *Select) buildPartsStr() string {
	builder := strings.Builder{}

	for _, part := range partsSort {
		partStr := sel.parts[part]
		if partStr != "" {
			if builder.Len() != 0 {
				builder.WriteRune(' ')
			}
			builder.WriteString(string(part))
			builder.WriteRune(' ')
			builder.WriteString(partStr)
		}
	}
	return builder.String()
}

func (sel *Select) concateStr(key SelectPart, query, join string) {
	str := sel.parts[key]
	if str == "" {
		str = query
	} else {
		str = str + join + query
	}
	sel.parts[key] = str
}
