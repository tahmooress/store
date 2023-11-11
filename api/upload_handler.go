package api

import (
	"bufio"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/tahmooress/store/internal/model"
	"github.com/tahmooress/store/internal/service"
)

func (h *HTTPServer) UploadFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerID, ok := getOwnerID(w, r)
		if !ok {
			return
		}

		err := r.ParseMultipartForm(int64(h.fileSizeLimit))
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("file size exceed the limit"))
			return
		}

		name := r.FormValue("name")
		if name == "" {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("name field of form is empty"))
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte("error retrieving the File, missing \"file\" key/value"))
			return
		}

		defer file.Close()

		fileObj := model.FileObject{
			OwnerID:   ownerID,
			Name:      name,
			Type:      r.FormValue("type"),
			Tags:      strings.Split(r.FormValue("tags"), ","),
			Content:   bufio.NewReader(file),
			Size:      handler.Size,
			CreatedAt: time.Now(),
		}

		if err := h.usecase.StoreFileObject(r.Context(), &fileObj); err != nil {
			if errors.Is(err, service.ErrInsufficientStorage) {
				w.WriteHeader(http.StatusInsufficientStorage)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("something goes wrong, try again!"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("success"))
	}
}

func getOwnerID(w http.ResponseWriter, r *http.Request) (string, bool) {
	ownerID := r.Header.Get("owner_id")
	if ownerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("owner_id header is not set"))
		return "", false
	}

	return ownerID, true
}
