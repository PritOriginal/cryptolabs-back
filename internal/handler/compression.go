package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/PritOriginal/cryptolabs-back/internal/services"
	"github.com/PritOriginal/problem-map-server/pkg/handlers"
	"github.com/PritOriginal/problem-map-server/pkg/responses"
)

type CompressionHandler struct {
	handlers.BaseHandler
	s services.Compression
}

func NewCompressionHandler(log *slog.Logger, s services.Compression) *CompressionHandler {
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

		dataCompressed := h.s.Compress(data)

		h.Render(w, r, responses.SucceededRenderer(string(dataCompressed)))
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

		h.Render(w, r, responses.SucceededRenderer(string(data)))
	}
}
