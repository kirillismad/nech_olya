package notesdb

import (
	"context"
	"database/sql"
	"fmt"
)

// Реализовать функцию SlowCount(ctx context.Context, db *sql.DB) (int64, error):
// - выполняет искусственно долгий запрос через рекурсивный CTE, например:
// sql
// WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM r WHERE n < 10000000)
// SELECT count(*) FROM r
// - использует QueryRowContext и возвращает результат Scan
// - корректно пробрасывает ошибку контекста

func SlowCount(ctx context.Context, db *sql.DB) (int64, error) {
	const query = `WITH RECURSIVE r(n) AS (SELECT 1 UNION ALL SELECT n+1 FROM r WHERE n < 10000000) SELECT count(*) FROM r`
	var count int64
	err := db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query failed: %w", err)
	}
	return count, nil
}
