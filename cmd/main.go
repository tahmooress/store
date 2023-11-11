package main

import (
	"context"
	"log"
	"os"

	"github.com/tahmooress/store/api"
	"github.com/tahmooress/store/internal/logger"
	"github.com/tahmooress/store/internal/repository/psql"
	"github.com/tahmooress/store/internal/service"
)

func main() {
	cfg := psql.Config{
		DatabaseName:        os.Getenv("DATABASE_NAME"),
		DatabaseUser:        os.Getenv("DATABASE_USER"),
		DatabasePassword:    os.Getenv("DATABASE_PASS"),
		DatabaseHost:        os.Getenv("DATABASE_HOST"),
		DatabasePort:        os.Getenv("DATABASE_PORT"),
		DatabaseMaxPageSize: os.Getenv("DATABASE_MAX_PAGE_SIZE"),
		DatabaseSSLMode:     os.Getenv("DATABASE_SSLMODE"),
	}

	repo, err := psql.New(context.Background(), &cfg)
	if err != nil {
		log.Fatal(err)
	}

	storageLimit := os.Getenv("STORAGE_LIMIT")
	dir := os.Getenv("STORAGE_DIR")
	pass := os.Getenv("FILE_ENCRYPTION_PASSWORD")
	logPath := os.Getenv("APP_LOG_PATH")
	logLevel := os.Getenv("LOG_LEVEL")

	logger, err := logger.New(logger.Config{
		LogFilePath: logPath,
		LogLevel:    logLevel,
	})
	if err != nil {
		log.Fatal(err)
	}

	usecase, err := service.New(repo, storageLimit, dir, []byte(pass), logger)
	if err != nil {
		log.Fatal(err)
	}

	ip := os.Getenv("HTTP_IP")
	port := os.Getenv("HTTP_PORT")
	uploadFileLimitSize := os.Getenv("UPLOAD_FILE_LIMIT_SIZE")

	httpServer, err := api.NewHTTPServer(usecase, logger, api.WithIP(ip), api.WithPort(port),
		api.WithUploadFileSizeLimit(uploadFileLimitSize))
	if err != nil {
		log.Fatal(err)
	}

	httpServer.ListenAndServe()

	err = <-httpServer.Err()
	if err != nil {
		log.Fatal(err)
	}
}
