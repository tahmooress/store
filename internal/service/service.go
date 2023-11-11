package service

import (
	"fmt"
	"os"

	"github.com/inhies/go-bytesize"
	"github.com/tahmooress/store/internal/logger"
	"github.com/tahmooress/store/internal/repository"
)

type Service struct {
	repo           repository.FileObjectStorage
	storageLimit   bytesize.ByteSize
	dir            string
	encryptionPass []byte
	logger         logger.Logger
}

func New(repo repository.FileObjectStorage, totalStorage, dir string, pass []byte,
	logger logger.Logger) (*Service, error) {
	limit, err := bytesize.Parse(totalStorage)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("service New parseLimitStorage: %s", err)
	}

	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			logger.Error(err)
			return nil, fmt.Errorf("service New parseLimitStorage: %s", err)
		}
	} else if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("service New parseLimitStorage: %s", err)
	}

	return &Service{
		repo:           repo,
		storageLimit:   limit,
		dir:            dir,
		encryptionPass: pass,
		logger:         logger,
	}, nil
}
