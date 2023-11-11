package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/inhies/go-bytesize"
	"github.com/tahmooress/store/internal/model"
	"github.com/tahmooress/store/internal/repository"
)

func (s *Service) FetchFileObject(ctx context.Context, filter *model.Filter) (
	[]model.FileObject, error,
) {
	fileObjs, err := s.repo.FetchFileObject(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("service FetchFileObject: %s", err)
	}

	for i := range fileObjs {
		fo, err := s.fsLoad(&fileObjs[i])
		if err != nil {
			return nil, fmt.Errorf("service FetchFileObject: %s", err)
		}

		fileObjs[i] = *fo
	}

	return fileObjs, nil
}

func (s *Service) StoreFileObject(ctx context.Context, fileObj *model.FileObject) error {
	fileObj.StorageLocation = filepath.Join(s.dir, fmt.Sprintf("%s.%s", fileObj.Name, fileObj.Type))
	err := s.repo.Exec(ctx, func(tx repository.Tx) error {
		size, err := tx.TotalStorage()
		if err != nil {
			return err
		}

		if s.storageLimit <= bytesize.ByteSize(size) {
			return ErrInsufficientStorage
		}

		if err := s.fsStore(fileObj); err != nil {
			return err
		}

		if err := tx.UpdateTotalStorage(fileObj.Size); err != nil {
			return err
		}

		return tx.StoreFileObject(fileObj)
	})
	if err != nil {
		return fmt.Errorf("service StoreFileObject :%s", err)
	}

	return nil
}

func (s *Service) fsStore(fileObject *model.FileObject) error {
	f, err := os.OpenFile(fileObject.StorageLocation, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("fsStore: %s", err)
	}

	defer f.Close()

	content, err := io.ReadAll(fileObject.Content)
	if err != nil {
		return fmt.Errorf("fsStore: %s", err)
	}

	encContenct, err := encryptWithAES(s.encryptionPass, content)
	if err != nil {
		return fmt.Errorf("fsStore: %s", err)
	}

	if _, err := f.Write(encContenct); err != nil {
		return fmt.Errorf("fsStore: %s", err)
	}
	return nil
}

func (s *Service) fsLoad(fileObject *model.FileObject) (*model.FileObject, error) {
	content, err := readFile(fileObject.StorageLocation)
	if err != nil {
		return nil, fmt.Errorf("fsLoad: %s", err)
	}

	b, err := decryptAES(content, s.encryptionPass)
	if err != nil {
		return nil, fmt.Errorf("fsLoad: %s", err)
	}

	fo := *fileObject
	fo.Content = bytes.NewReader(b)

	return &fo, nil
}

func (s *Service) DeleteFileObject(ctx context.Context, fileObj []model.FileObject) error {
	var idx []int64

	for _, fo := range fileObj {
		idx = append(idx, fo.ID)
	}

	err := s.repo.Exec(ctx, func(tx repository.Tx) error {
		size, err := tx.BatchDeleteFileObjects(idx...)
		if err != nil {
			return err
		}

		for _, fo := range fileObj {
			if err := os.RemoveAll(fo.StorageLocation); err != nil {
				return err
			}
		}

		return tx.UpdateTotalStorage(-size)
	})
	if err != nil {
		return fmt.Errorf("service DeleteFileObject :%s", err)
	}

	return nil
}

func readFile(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("createFileObject: %s", err)
	}

	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("createFileObject: %s", err)
	}

	return content, nil
}
