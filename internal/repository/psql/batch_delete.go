package psql

import (
	"fmt"
)

func (t *transaction) BatchDeleteFileObjects(idx ...int64) (int64, error) {
	qb := newQueryBuilder("DELETE FROM file_object ")

	var args []interface{}

	for _, id := range idx {
		args = append(args, id)
	}

	qb.wherein("file_object.id", args)
	qb.appendQuery(" RETURNING size")

	query, args, err := qb.build()
	if err != nil {
		return 0, fmt.Errorf("psql BatchDeleteFileObjects beginTx :%s", err)
	}

	rows, err := t.tx.Query(query, args...)
	if err != nil {
		return 0, fmt.Errorf("psql BatchDeleteFileObjects Query :%s", err)
	}

	defer rows.Close()

	var freedSpace int64

	for rows.Next() {
		var fs int64

		if err := rows.Scan(&fs); err != nil {
			return 0, fmt.Errorf("psql BatchDeleteFileObjects Scan :%s", err)
		}

		freedSpace += fs
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("psql BatchDeleteFileObjects rows.Err :%s", err)
	}

	qb.reset()

	qb.appendQuery("DELETE FROM file_tag ")
	qb.wherein("file_tag.file_object_id", args)

	q, arg, err := qb.build()
	if err != nil {
		return 0, fmt.Errorf("psql BatchDeleteFileObjects beginTx :%s", err)
	}

	if _, err = t.tx.Exec(q, arg...); err != nil {
		return 0, fmt.Errorf("psql BatchDeleteFileObjects beginTx :%s", err)
	}

	return freedSpace, nil
}
