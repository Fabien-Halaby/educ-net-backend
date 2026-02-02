package repository

import (
	"database/sql"
	"fmt"
)

// ! domainScanner interface générique pour tous les repositories
type domainScanner interface {
	Scan(dest ...any) error
}

// ! nullString convertit sql.NullString → string (empty si NULL)
func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// ! nullInt convertit sql.NullInt64 → *int (nil si NULL)
func nullInt(ni sql.NullInt64) *int {
	if ni.Valid {
		v := int(ni.Int64)
		return &v
	}
	return nil
}

// ! scanError wrapper standard pour tous les scan
func scanError(err error, operation string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}
