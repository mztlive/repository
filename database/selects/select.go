package selects

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/mztlive/repository/database"
)

// ByIdentities 根据唯一编号查询某个Domain List
// 如何是聚合对象不要用这个方法，因为这个方法是查询单个Domain类型的
func ByIdentities[T interface{}](ctx context.Context, identities []string, tableName string) ([]*T, error) {

	result := []*T{}
	sql, args, err := squirrel.Select("*").From(tableName).Where(squirrel.Eq{"identity": identities}).ToSql()

	if err != nil {
		return result, fmt.Errorf("generator SQL failed, %w", err)
	}

	err = database.GetDB().SelectContext(ctx, &result, sql, args...)
	return result, err
}

// CollectionByColumn 根据某个字段查询多行数据
func CollectionByColumn[T any](ctx context.Context, val, column, tableName string) ([]*T, error) {
	result := make([]*T, 0)
	runSQL := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?;", tableName, column)
	err := database.GetDB().SelectContext(ctx, &result, runSQL, val)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

// ByIdentity 根据Identity查询
func ByIdentity[T interface{}](ctx context.Context, identity string, tableName string) (*T, error) {
	result := new(T)
	rawSql := "SELECT * FROM " + tableName + " WHERE identity = ?;"
	err := database.GetDB().GetContext(
		ctx,
		result,
		rawSql,
		identity,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

// ByID 根据ID查询
func ByID[T interface{}](ctx context.Context, id int64, tableName string) (*T, error) {
	result := new(T)

	err := database.GetDB().GetContext(
		ctx,
		result,
		"SELECT * FROM "+tableName+" WHERE id = ?;",
		id,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}
