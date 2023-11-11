package psql

import (
	"database/sql"
	"errors"
	"fmt"
)

func (t *transaction) TotalStorage() (int64, error) {
	query := "SELECT size FROM total_storage"

	var size int64

	if err := t.tx.QueryRow(query).Scan(&size); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("psql TotalStorage Scan: %s", err)
	}

	return size, nil
}

func (t *transaction) UpdateTotalStorage(diff int64) error {
	query := "UPDATE total_storage SET size = size + $1"

	if _, err := t.tx.Exec(query, diff); err != nil {
		return fmt.Errorf("psql UpdateTotalStorage: %s", err)
	}

	return nil
}
