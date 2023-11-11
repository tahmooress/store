package service

import (
	"context"
	"errors"

	"github.com/tahmooress/store/internal/model"
)

var ErrInsufficientStorage = errors.New("insufficient storage")

type Usecase interface {
	FetchFileObject(ctx context.Context, filter *model.Filter) ([]model.FileObject, error)
	StoreFileObject(ctx context.Context, fileObj *model.FileObject) error
	DeleteFileObject(ctx context.Context, fileObj []model.FileObject) error
}
