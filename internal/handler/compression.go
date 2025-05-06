package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/PritOriginal/cryptolabs-back/internal/services/compression"
	"github.com/PritOriginal/problem-map-server/pkg/handlers"
	"github.com/PritOriginal/problem-map-server/pkg/responses"
)

type CompressionService interface {
	Compress(data []byte) ([]byte, error)
	CompressWithDetails(data []byte) (compression.CompressionDetails, error)
	Decompress(compressedData []byte) ([]byte, error)
	DecompressWithDetails(data []byte) (compression.CompressionDetails, error)
}

type CompressionHandler struct {
	handlers.BaseHandler
	s CompressionService
}

func NewCompressionHandler(log *slog.Logger, s CompressionService) *CompressionHandler {
	return &CompressionHandler{handlers.BaseHandler{Log: log}, s}
}

func (h *CompressionHandler) Compress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			h.RenderError(w, r,
				handlers.HandlerError{Msg: "invalid data", Err: err},
				responses.ErrBadRequest,
			)
			return
		}

		dataCompressed, err := h.s.Compress(data)
		if err != nil {
			h.RenderInternalError(w, r, handlers.HandlerError{Msg: "failed compress", Err: err})
			return
		}

		h.Log.Debug("size compress", slog.Int("size", len(dataCompressed)))

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=test.txt")
		w.Write(dataCompressed)
	}
}

func (h *CompressionHandler) CompressWithDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			h.RenderError(w, r,
				handlers.HandlerError{Msg: "invalid data", Err: err},
				responses.ErrBadRequest,
			)
			return
		}

		details, err := h.s.CompressWithDetails(data)
		if err != nil {
			h.RenderInternalError(w, r, handlers.HandlerError{Msg: "failed compress", Err: err})
			return
		}

		h.Log.Debug("size compress", slog.Int("size", len(details.Data)))
		h.renderDetails(w, r, details)
	}
}

func (h *CompressionHandler) Decompress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataCompressed, err := io.ReadAll(r.Body)
		if err != nil {
			h.RenderError(w, r,
				handlers.HandlerError{Msg: "invalid data", Err: err},
				responses.ErrBadRequest,
			)
			return
		}

		data, err := h.s.Decompress(dataCompressed)
		if err != nil {
			h.RenderError(w, r,
				handlers.HandlerError{Msg: "invalid data", Err: err},
				responses.ErrBadRequest,
			)
			return
		}
		h.Log.Debug("size decompress", slog.Int("size", len(data)))

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Disposition", "attachment; filename=test.txt")
		w.Write(data)
	}
}

func (h *CompressionHandler) DecompressWithDetails() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataCompressed, err := io.ReadAll(r.Body)
		if err != nil {
			h.RenderError(w, r,
				handlers.HandlerError{Msg: "invalid data", Err: err},
				responses.ErrBadRequest,
			)
			return
		}

		details, err := h.s.DecompressWithDetails(dataCompressed)
		if err != nil {
			h.RenderError(w, r,
				handlers.HandlerError{Msg: "invalid data", Err: err},
				responses.ErrBadRequest,
			)
			return
		}
		h.Log.Debug("size decompress", slog.Int("size", len(details.Data)))
		h.renderDetails(w, r, details)
	}
}

func (h *CompressionHandler) renderDetails(w http.ResponseWriter, r *http.Request, details compression.CompressionDetails) {
	mpw := multipart.NewWriter(w)
	defer mpw.Close()
	w.Header().Set("Content-Type", mpw.FormDataContentType())

	detailsJson, err := json.Marshal(details.Details)
	if err != nil {
		h.RenderInternalError(w, r, handlers.HandlerError{Msg: "failed marshal details", Err: err})
		return
	}
	mpw.WriteField("details", string(detailsJson))

	fWriter, err := mpw.CreateFormFile("data", "test.txt")
	if err != nil {
		h.RenderInternalError(w, r, handlers.HandlerError{Msg: "failed create form", Err: err})
		return
	}

	fBuf := bytes.NewBuffer(details.Data)
	_, err = io.Copy(fWriter, fBuf)
	if err != nil {
		h.RenderInternalError(w, r, handlers.HandlerError{Msg: "failed write data", Err: err})
		return
	}
}
