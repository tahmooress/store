package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/tahmooress/store/internal/model"
	"github.com/tahmooress/store/internal/repository"
)

func (d *DB) FetchFileObject(ctx context.Context, filter *model.Filter) (
	[]model.FileObject, error,
) {
	wqStr, args, err := filter.Query(fetchObjectfilterQuery)
	if err != nil {
		return nil, fmt.Errorf("psql FetchAndDeleteFileObject: %s", err)
	}

	rows, err := d.conn.QueryContext(ctx, wqStr, args...)
	if err != nil {
		return nil, fmt.Errorf("psql FetchAndDeleteFileObject: %s", err)
	}

	defer rows.Close()

	var fileObjects []model.FileObject

	for rows.Next() {
		var fo model.FileObject

		err := rows.Scan(&fo.ID, &fo.OwnerID, &fo.Name, &fo.Type, &fo.Size,
			&fo.StorageLocation, &fo.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("psql FetchAndDeleteFileObject: %s", err)
		}

		fileObjects = append(fileObjects, fo)
	}

	if err := rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("psql FetchAndDeleteFileObject: %s", repository.ErrNotFound)
		}

		return nil, fmt.Errorf("psql FetchAndDeleteFileObject: %s", err)
	}

	return fileObjects, nil
}

func fetchObjectfilterQuery(owner_id, name string, tags []string) (
	string, []interface{}, error,
) {
	baseQuery := `
		SELECT file_object.id, owner_id, file_object.name, type, size,storage_location,
		created_at FROM
		file_object INNER JOIN  
		file_tag ft ON  file_object.id = ft.file_object_id 
		INNER JOIN tag ON ft.tag_id = tag.id 
	`

	qb := newQueryBuilder(baseQuery)

	qb.whereCond("owner_id", owner_id)

	if name != "" || len(tags) > 0 {
		qb.and()
	}

	if name != "" {
		if len(tags) > 0 {
			qb.stratParentheses()
		}

		qb.whereCond("file_object.name", name)
	}

	if len(tags) != 0 {
		if name != "" {
			qb.or()
		}

		var tgs []interface{}

		for _, t := range tags {
			tgs = append(tgs, t)
		}

		qb.wherein("tag.name", tgs)

		if name != "" {
			qb.closeParentheses()
		}
	}

	return qb.build()
}
