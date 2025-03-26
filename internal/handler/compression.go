package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/PritOriginal/cryptolabs-back/internal/services"
	"github.com/PritOriginal/problem-map-server/pkg/logger"
	"github.com/PritOriginal/problem-map-server/pkg/responses"
	"github.com/go-chi/render"
)

type CompressionHandler struct {
	log *slog.Logger
	s   services.Compression
}

func NewCompressionHandler(log *slog.Logger, s services.Compression) *CompressionHandler {
	return &CompressionHandler{log, s}
}

func (h *CompressionHandler) render(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	if err := render.Render(w, r, v); err != nil {
		h.log.Error("failed render", logger.Err(err))
		render.Render(w, r, responses.ErrInternalServer)
	}
}

func (h *CompressionHandler) renderBadRequest(w http.ResponseWriter, r *http.Request, msg string, err error) {
	h.log.Error(msg, logger.Err(err))
	render.Render(w, r, responses.ErrBadRequest)
}

func (h *CompressionHandler) Compress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			h.renderBadRequest(w, r, "invalid data", err)
			return
		}

		dataCompressed := h.s.Compress(data)

		h.render(w, r, responses.SucceededRenderer(string(dataCompressed)))
	}
}

func (h *CompressionHandler) Decompress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dataCompressed, err := io.ReadAll(r.Body)
		if err != nil {
			h.renderBadRequest(w, r, "invalid data", err)
			return
		}

		data, err := h.s.Decompress(dataCompressed)
		if err != nil {
			h.renderBadRequest(w, r, "invalid data", err)
			return
		}

		h.render(w, r, responses.SucceededRenderer(string(data)))
	}
}
