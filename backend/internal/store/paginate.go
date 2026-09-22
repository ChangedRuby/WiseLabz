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
	total, err := countRows(ctx, db, table, where, args)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT " + columns + " FROM " + table + " " + where + " ORDER BY " + orderBy + " LIMIT ? OFFSET ?"
	items, err := scanAll(ctx, db, table, query, append(args, limit, offset), scan)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Keyset is a cursor position for keyset pagination: the sort-column value
// and row id of the last row the client already saw. A zero Keyset means
// "start at the newest row".
type Keyset struct {
	Sort string
	ID   string
}

// Empty reports whether the keyset is the start-of-stream position.
func (k Keyset) Empty() bool { return k.Sort == "" || k.ID == "" }

// keysetQuery is the paginatedQuery shape for keyset (cursor) pagination:
// SELECT columns FROM table <where> [AND (sortCol, id) < (?, ?)]
// ORDER BY sortCol DESC, id DESC LIMIT ?.
//
// The row-value comparison is what makes the page boundary exact when several
// rows share a sortCol value; it is supported by both SQLite (3.15+) and
// PostgreSQL. total is the count of all rows matching where, ignoring the
// cursor, so the envelope keeps reporting the full result size.
func keysetQuery[T any](ctx context.Context, db DBTX, table, columns, where string, args []any, sortCol string, cur Keyset, limit int, scan func(rowScanner) (T, error)) ([]T, int, error) {
	total, err := countRows(ctx, db, table, where, args)
	if err != nil {
		return nil, 0, err
	}

	pageWhere, pageArgs := where, args
	if !cur.Empty() {
		if pageWhere == "" {
			pageWhere = "WHERE"
		} else {
			pageWhere += " AND"
		}
		pageWhere += " (" + sortCol + ", id) < (?, ?)"
		pageArgs = append(append([]any{}, args...), cur.Sort, cur.ID)
	}

	query := "SELECT " + columns + " FROM " + table + " " + pageWhere +
		" ORDER BY " + sortCol + " DESC, id DESC LIMIT ?"
	items, err := scanAll(ctx, db, table, query, append(pageArgs, limit), scan)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func countRows(ctx context.Context, db DBTX, table, where string, args []any) (int, error) {
	var total int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+table+" "+where, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count %s: %w", table, err)
	}
	return total, nil
}

// scanAll runs query and scans every row with scan. Always returns a non-nil,
// possibly empty slice.
func scanAll[T any](ctx context.Context, db DBTX, table, query string, args []any, scan func(rowScanner) (T, error)) ([]T, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list %s: %w", table, err)
	}
	defer rows.Close() //nolint:errcheck

	items := make([]T, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("scan %s: %w", table, err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate %s: %w", table, err)
	}
	return items, nil
}
