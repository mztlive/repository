package database

import (
	"fmt"
	"reflect"
	"strings"
)

// BuilderInsertSQL 通过反射生成插入语句
// 传入的参数必须是struct
// 返回的是一个sql, args, error
func BuilderInsertSQL(ptr interface{}) (string, []interface{}, error) {
	reType := reflect.TypeOf(ptr)
	if reType.Kind() != reflect.Ptr || reType.Elem().Kind() != reflect.Struct {
		return "", nil, fmt.Errorf("obj must be a struct pointer")
	}

	tableName := camelToUnderline(reType.Elem().Name()) + "s"
	colBuf := []string{}
	valueBuf := []interface{}{}

	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		// 获取结构体字段信息
		structField := v.Type().Field(i)
		// 取tag
		tag := structField.Tag

		dbField := tag.Get("db")
		if dbField == "" || dbField == "-" {
			continue
		}

		structFieldValue := v.Field(i)

		switch structField.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if structFieldValue.Int() == 0 && dbField == "id" {
				continue
			}
			valueBuf = append(valueBuf, structFieldValue.Int())
		case reflect.String:
			valueBuf = append(valueBuf, structFieldValue.String())
		case reflect.Float32, reflect.Float64:
			valueBuf = append(valueBuf, structFieldValue.Float())
		case reflect.Bool:
			if structFieldValue.Bool() {
				valueBuf = append(valueBuf, 1)
			} else {
				valueBuf = append(valueBuf, 0)
			}
		case reflect.Ptr:
			if structFieldValue.IsNil() {
				continue
			}
			valueBuf = append(valueBuf, v.Field(i).Interface())
		}

		colBuf = append(colBuf, dbField)
	}

	if len(colBuf) == 0 && len(valueBuf) == 0 {
		return "", nil, fmt.Errorf("No Builder. Check your struct")
	}

	placeholder := strings.Repeat("?,", len(colBuf))
	placeholder = placeholder[:len(placeholder)-1]

	sql := "insert into " + tableName + "(" + strings.Join(colBuf, ",") + ") values(" + placeholder + ")"
	return sql, valueBuf, nil
}

// camelToUnderline 驼峰转下划线
func camelToUnderline(s string) string {
	var buf []byte
	for i := 0; i < len(s); i++ {
		if s[i] >= 'A' && s[i] <= 'Z' {
			if i > 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, s[i])
		} else {
			buf = append(buf, s[i])
		}
	}
	return strings.ToLower(string(buf))
}
