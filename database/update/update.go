package update

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/mztlive/repository/database"
)

// UpdateSingleFieldPayload 更新单个字段的Payload
type UpdateSingleFieldPayload[T any] struct {
	Identities      []string
	IdentityColName string
	Table           string
	FieldName       string
	FieldValue      T
}

// UpdateFieldByIdentities 更新单个字段
func UpdateFieldByIdentities[T any](ctx context.Context, payload UpdateSingleFieldPayload[T]) error {
	sql, args, err := squirrel.Update(payload.Table).
		Set(payload.FieldName, payload.FieldValue).
		Where(squirrel.Eq{payload.IdentityColName: payload.Identities}).
		ToSql()

	if err != nil {
		return fmt.Errorf("generator SQL failed, %w", err)
	}

	result, err := database.GetDB().ExecContext(ctx, sql, args...)
	return CheckErr(result, err, payload.Table)
}

// CheckErr 封装了更新结果检查的逻辑， 只适用于单个domain的更新结果的检查
func CheckErr(result sql.Result, err error, tableName string) error {
	if err != nil {
		return fmt.Errorf("update %s failed: %w", tableName, err)
	}

	// affected, _ := result.RowsAffected()
	// if affected == 0 {
	// 	return fmt.Errorf("Update %s Failed: No Rows Affected", tableName)
	// }

	return nil
}
