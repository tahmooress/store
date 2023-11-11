package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/inhies/go-bytesize"
	"github.com/tahmooress/store/internal/logger"
	"github.com/tahmooress/store/internal/service"
)

const (
	defaultReadTimeout     = 30 * time.Second
	defaultWriteTimeout    = 120 * time.Second
	defaultIP              = "localhost"
	defaultPort            = "8080"
	defaultUploadFilelimit = bytesize.GB
)

type HTTPServer struct {
	ip            string
	port          string
	readTimeout   time.Duration
	writeTimeout  time.Duration
	httpSrv       *http.Server
	router        *mux.Router
	usecase       service.Usecase
	errc          chan error
	fileSizeLimit bytesize.ByteSize
	logger        logger.Logger
}

func NewHTTPServer(service service.Usecase, logger logger.Logger, ops ...Option) (
	*HTTPServer, error,
) {
	server := &HTTPServer{
		usecase:       service,
		ip:            defaultIP,
		port:          defaultPort,
		readTimeout:   defaultReadTimeout,
		writeTimeout:  defaultWriteTimeout,
		router:        mux.NewRouter(),
		errc:          make(chan error),
		fileSizeLimit: defaultUploadFilelimit,
		logger:        logger,
	}

	for _, o := range ops {
		o(server)
	}

	server.router.HandleFunc("/upload", server.UploadFileHandler()).Methods(http.MethodPost)
	server.router.HandleFunc("/download", server.FetchFileHandler()).Methods(http.MethodGet)

	return server, nil
}

func (h *HTTPServer) Close() error {
	if h.httpSrv == nil {
		return nil
	}
	return h.httpSrv.Close()
}

func (h *HTTPServer) Err() <-chan error {
	return h.errc
}

func (h *HTTPServer) ListenAndServe() {
	h.httpSrv = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", h.ip, h.port),
		Handler:      h.router,
		ReadTimeout:  h.readTimeout,
		WriteTimeout: h.writeTimeout,
	}

	go func() {
		h.errc <- h.httpSrv.ListenAndServe()
	}()

	h.logger.Printf("http server run on ip: %s port: %s", h.ip, h.port)
}

type Option func(h *HTTPServer)

func WithPort(port string) Option {
	return func(h *HTTPServer) {
		h.port = port
	}
}

func WithIP(ip string) Option {
	return func(h *HTTPServer) {
		h.ip = ip
	}
}

func WithUploadFileSizeLimit(size string) Option {
	return func(h *HTTPServer) {
		s, err := bytesize.Parse(strings.ToUpper(size))
		if err == nil {
			h.fileSizeLimit = s
		}
	}
}
