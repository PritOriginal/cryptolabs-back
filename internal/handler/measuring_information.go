package handler

import (
	"log/slog"
	"net/http"

	"github.com/PritOriginal/cryptolabs-back/internal/services"
	"github.com/PritOriginal/problem-map-server/pkg/logger"
	"github.com/PritOriginal/problem-map-server/pkg/responses"
	"github.com/go-chi/render"
)

type MeasuringInformationHandler struct {
	log *slog.Logger
	uc  services.MeasuringInformation
}

func NewMeasuringInformation(log *slog.Logger, uc services.MeasuringInformation) *MeasuringInformationHandler {
	return &MeasuringInformationHandler{log, uc}
}

func (h *MeasuringInformationHandler) render(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	if err := render.Render(w, r, v); err != nil {
		h.log.Error("failed render", logger.Err(err))
		render.Render(w, r, responses.ErrInternalServer)
	}
}

func (h *MeasuringInformationHandler) GetAlphabet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alphabetSet := r.URL.Query().Get("alphabet_set")

		alphabet, err := h.uc.GetAlphabet(alphabetSet, "")
		if err != nil {
			return
		}

		h.render(w, r, responses.SucceededRenderer(alphabet))
	}
}

func (h *MeasuringInformationHandler) GetInformationVolumeSymbol() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alphabetSet := r.URL.Query().Get("alphabet_set")
		customAlphabet := r.URL.Query().Get("alphabet")

		alphabet, err := h.uc.GetAlphabet(alphabetSet, customAlphabet)
		if err != nil {
			h.log.Error("error get alphabet", logger.Err(err))
			render.Render(w, r, responses.ErrorRenderer(err))
			return
		}

		volume := h.uc.GetInformationVolumeSymbol(alphabet)

		h.render(w, r, responses.SucceededRenderer(volume))
	}
}

func (h *MeasuringInformationHandler) GetAmountOfInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query().Get("text")
		alphabetSet := r.URL.Query().Get("alphabet_set")
		alphabet_param := r.URL.Query().Get("alphabet")

		alphabet, err := h.uc.GetAlphabet(alphabetSet, alphabet_param)
		if err != nil {
			h.log.Error("error get alphabet", logger.Err(err))
			render.Render(w, r, responses.ErrorRenderer(err))
			return
		}

		amount := h.uc.GetAmountOfInformation(text, alphabet)
		h.render(w, r, responses.SucceededRenderer(amount))
	}
}
