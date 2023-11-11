package api

import (
	"archive/zip"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tahmooress/store/internal/model"
)

func (h *HTTPServer) FetchFileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerID, ok := getOwnerID(w, r)
		if !ok {
			return
		}
		filter := model.NewFilter(ownerID).
			SetName(r.URL.Query().Get("name")).
			SetTags(strings.Split(r.URL.Query().Get("tags"), ","))

		fileObjs, err := h.usecase.FetchFileObject(r.Context(), filter)
		if err != nil {
			h.logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("something goes wrong, try again"))
			return
		}
		if fileObjs == nil {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte("not found"))
			return
		}

		d := url.PathEscape("download.zip")
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename="+d+";")

		z := zip.NewWriter(w)
		defer z.Close()

		for _, fo := range fileObjs {
			zf, err := z.Create(fo.Name)
			if err != nil {
				h.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("something goes wrong, try again"))
				return
			}

			_, err = io.Copy(zf, fo.Content)
			if err != nil {
				h.logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte("something goes wrong, try again"))
				return
			}
		}

		go func() {
			err := h.usecase.DeleteFileObject(context.TODO(), fileObjs)
			if err != nil {
				h.logger.Error(err)
			}
		}()
	}
}
