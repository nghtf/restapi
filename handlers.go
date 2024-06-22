package restapi

import (
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type TGenericHandlers struct {
	api *TRestAPI
}

/*
	Generic GET handler
*/

// Generic GET handler. Returns data to a client (automatically marshalled as TResponseTemplate.Data)
// Ex1: {"status":"Ok","data":"simple string"}
// Ex2: {"status":"Ok","data":{"field1":"value","field2":4,"field3":1.3}}
func (h *TGenericHandlers) GET(data interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := enrich(h.api.log, r)
		log.Info("new request")
		render.JSON(w, r, TResponse{}.Data(data))
	}
}

/*
	Generic POST handler
*/

const DEFAULT_UPLOAD_PATTERN = "upload_*"

type TFileUpload struct {
	ID     string // UUID
	Path   string // path to uploaded file
	Header *multipart.FileHeader
}

// Generic POST handler. Retrieves a file and (optionally) sends filepath to the channel.
// Uploaded files are stored according to name pattern provided (or default DEFAULT_UPLOAD_PATTERN).
func (h *TGenericHandlers) POST(formField string, uploadDir string, tempPattern string, fch chan TFileUpload) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log := enrich(h.api.log, r)
		log.Info("new request")

		fu, err := uploader(formField, uploadDir, tempPattern, r)
		if err != nil {
			render.JSON(w, r, TResponse{}.Error("internal error"))
			log.Error("upload failed", e(err))
			return
		}

		log.Info("new upload",
			slog.Group("file",
				slog.String("UUID", fu.ID),
				slog.String("Path", fu.Path),
				slog.String("Filename", fu.Header.Filename),
				slog.Int64("Size", fu.Header.Size),
			),
		)

		if fch != nil {
			fch <- fu
		}

		render.JSON(w, r, TResponse{}.Data(fu.ID))
	}
}

func uploader(formField string, uploadDir string, tempPattern string, r *http.Request) (TFileUpload, error) {

	var fileUpload TFileUpload

	r.ParseMultipartForm(10 << 20)

	if _, ok := r.MultipartForm.File[formField]; !ok {
		return fileUpload, errors.New("form field not found")
	}

	file, header, err := r.FormFile(formField)
	if err != nil {
		return fileUpload, err
	}
	defer file.Close()

	fileUpload.Header = header

	if tempPattern == "" {
		tempPattern = DEFAULT_UPLOAD_PATTERN
	}

	dst, err := os.CreateTemp(uploadDir, tempPattern)
	if err != nil {
		return fileUpload, err
	}
	defer dst.Close()

	fileUpload.Path = dst.Name()

	if _, err := io.Copy(dst, file); err != nil {
		return fileUpload, err
	}

	fileUpload.ID = uuid.New().String()

	return fileUpload, nil
}

// tools

func e(err error) slog.Attr {
	return slog.String("error", err.Error())
}

func enrich(log *slog.Logger, r *http.Request) *slog.Logger {
	log = log.With(
		slog.String("handler", "generic"),
		slog.Group("request",
			slog.String("id", middleware.GetReqID(r.Context())),
			slog.String("route", r.URL.Path),
			slog.String("method", r.Method),
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("user_agent", r.UserAgent()),
		),
	)
	return log
}
