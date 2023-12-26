package util

import (
	"fmt"
)

func AppendWhereTerm(term string, inQueryString string, inQueryConditions []interface{}, inIsFirst bool, value interface{}) (queryString string, queryConditions []interface{}, isFirst bool) {
	if inIsFirst {
		queryString = fmt.Sprintf("WHERE %s", term)
	} else {
		queryString = fmt.Sprintf("%s AND %s", inQueryString, term)
	}
	isFirst = false
	queryConditions = append(inQueryConditions, value)

	return
}

func AppendWhereSearchString(tableName string, inQueryString string, columns []string, inQueryConditions []interface{}, inIsFirst bool, searchString string) (keywordQuery string, queryConditions []interface{}, isFirst bool) {
	tableNameDot := ""
	if tableName != "" {
		tableNameDot = tableName + "."
	}

	queryConditions = inQueryConditions
	for i, column := range columns {
		if i > 0 {
			keywordQuery += " OR "
		}
		keywordQuery = fmt.Sprintf(`%s%s%s LIKE CONCAT("%%", ?, "%%")`, keywordQuery, tableNameDot, column)
		queryConditions = append(queryConditions, searchString)
	}
	if len(columns) > 1 {
		keywordQuery = fmt.Sprintf("(%s)", keywordQuery)
	}

	if inIsFirst {
		keywordQuery = fmt.Sprintf("WHERE %s", keywordQuery)
	} else {
		keywordQuery = fmt.Sprintf("%s AND %s", inQueryString, keywordQuery)
	}
	isFirst = false
	return
}
