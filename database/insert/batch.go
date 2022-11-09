package insert

import (
	"context"
	"fmt"

	"github.com/mztlive/repository/database"
)

// Batch 批量插入
//
// T 是任何domain。
//
// 通过事务执行多条插入语句完成的批量动作，这意味着每个domain的插入语句可以是不一样的
func Batch[T any](ctx context.Context, domains []*T) error {
	tx := database.GetDB().MustBegin()
	for _, domain := range domains {
		sql, args, err := database.BuilderInsertSQL(domain)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("generator SQL failed, %w", err)
		}

		if _, err = tx.ExecContext(ctx, sql, args...); err != nil {
			tx.Rollback()
			return fmt.Errorf("execute SQL failed, %w", err)
		}
	}
	return tx.Commit()
}
