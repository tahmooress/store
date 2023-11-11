package psql

import (
	"fmt"
	"strings"

	"github.com/tahmooress/store/internal/model"
)

func (t *transaction) StoreFileObject(fileObject *model.FileObject) error {
	tagsQuery := fmt.Sprintf("WITH input(name) AS (SELECT unnest('{%s}'::varchar(50)[]))",
		strings.Join(fileObject.Tags, ",")) +
		`, ins AS (
			INSERT INTO tag(name) 
			TABLE  input
			ON CONFLICT (name) DO NOTHING
			RETURNING id
			)
		 SELECT t.id
		 FROM   input i
		 JOIN   tag   t USING (name)
		 UNION  ALL
		 TABLE  ins
		`

	rows, err := t.tx.Query(tagsQuery)
	if err != nil {
		return fmt.Errorf("psql StoreFileObject inserting tags: %s", err)
	}

	defer rows.Close()

	var (
		ids []int64
		id  int64
	)

	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("psql StoreFileObject scan tag id: %s", err)
		}

		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("psql StoreFileObject rows.Err: %s", err)
	}

	foQuery := `INSERT INTO file_object(owner_id, name, type, size,storage_location)
				VALUES($1,$2,$3,$4,$5)	RETURNING id;	
	`
	var fileObjectID int64

	err = t.tx.QueryRow(foQuery, fileObject.OwnerID, fileObject.Name,
		fileObject.Type, fileObject.Size, fileObject.StorageLocation).Scan(&fileObjectID)
	if err != nil {
		return fmt.Errorf("psql StoreFileObject insert into file_object: %s", err)
	}

	var (
		indecies []string
		args     []interface{}
	)

	for i, id := range ids {
		indecies = append(indecies, fmt.Sprintf("($%d,$%d)", (i*2)+1, (i*2)+2))
		args = append(args, id, fileObjectID)
	}

	ftQuery := fmt.Sprintf("INSERT INTO file_tag(tag_id, file_object_id) VALUES %s",
		strings.Join(indecies, ","))

	_, err = t.tx.Exec(ftQuery, args...)
	if err != nil {
		return fmt.Errorf("psql StoreFileObject insert into tag_file: %s", err)
	}

	return nil
}
