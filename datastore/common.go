package datastore

import (
	"reflect"
	"strconv"
)

// makePrepareExpressionParamsForInOperand makes prepare expression for in operand
func makePrepareExpressionParamsForInOperand(target interface{}) (string, map[string]interface{}) {
	bindParams := make(map[string]interface{})
	query := ""
	rv := reflect.ValueOf(target)

	switch rv.Kind() {
	case reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if len(query) > 0 {
				query += ", "
			}
			bindParams["id"+strconv.Itoa(i)] = rv.Index(i)
			query += ":id" + strconv.Itoa(i)
		}
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			if len(query) > 0 {
				query += ", "
			}
			bindParams["id"+strconv.Itoa(i)] = rv.Index(i).Interface()
			query += ":id" + strconv.Itoa(i)
		}
	}
	return query, bindParams
}

// makePrepareExpressionForInOperand is make prepare for in expression
func makePrepareExpressionForInOperand(target interface{}) (string, []interface{}) {
	var bindParams []interface{}
	query := ""
	rv := reflect.ValueOf(target)

	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		bindParams = make([]interface{}, rv.Len(), rv.Len())
		for i := 0; i < rv.Len(); i++ {
			if len(query) > 0 {
				query += ", "
			}
			bindParams[i] = rv.Index(i).Interface()
			query += "?"
		}
	}
	return query, bindParams
}
