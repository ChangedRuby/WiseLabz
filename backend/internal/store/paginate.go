package store

import (
	"context"
	"fmt"
)

// paginatedQuery runs the count-then-select-then-scan shape shared by the
// store's paginated List* methods: SELECT COUNT(*) FROM table <where>, then
// SELECT columns FROM table <where> ORDER BY orderBy LIMIT ? OFFSET ?,
// scanning each row with scan. args are the WHERE clause's bind values (not
// including limit/offset); where may be "" for no filter. Always returns a
// non-nil, possibly empty slice.
func paginatedQuery[T any](ctx context.Context, db DBTX, table, columns, where string, args []any, orderBy string, limit, offset int, scan func(rowScanner) (T, error)) ([]T, int, error) {
	var total int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table+" "+where, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count %s: %w", table, err)
	}

	query := "SELECT " + columns + " FROM " + table + " " + where + " ORDER BY " + orderBy + " LIMIT ? OFFSET ?"
	rows, err := db.QueryContext(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, fmt.Errorf("list %s: %w", table, err)
	}
	defer rows.Close() //nolint:errcheck

	items := make([]T, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan %s: %w", table, err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate %s: %w", table, err)
	}
	return items, total, nil
}
