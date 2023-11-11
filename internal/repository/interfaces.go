package repository

import (
	"context"
	"errors"
	"io"

	"github.com/tahmooress/store/internal/model"
)

var ErrNotFound = errors.New("requested data not found")

type Tx interface {
	BatchDeleteFileObjects(idx ...int64) (int64, error)
	StoreFileObject(fileObj *model.FileObject) error
	TotalStorage() (int64, error)
	UpdateTotalStorage(diff int64) error
}

type FileObjectStorage interface {
	FetchFileObject(ctx context.Context, filter *model.Filter) ([]model.FileObject, error)
	Exec(ctx context.Context, fn func(tx Tx) error) error
	io.Closer
}
