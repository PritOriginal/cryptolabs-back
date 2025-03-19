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

// type MeasuringInformation

func (h *MeasuringInformationHandler) GetAlphabet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alphabetSet := r.URL.Query().Get("alphabet_set")

		alphabet, err := h.uc.GetAlphabet(alphabetSet)
		if err != nil {
			return
		}

		if err := render.Render(w, r, responses.SucceededRenderer(alphabet)); err != nil {
			h.log.Error("failed succeeded render", logger.Err(err))
			render.Render(w, r, responses.ErrInternalServer)
			return
		}
	}
}

func (h *MeasuringInformationHandler) GetInformationVolumeSymbol() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alphabetSet := r.URL.Query().Get("alphabet_set")
		alphabet_param := r.URL.Query().Get("alphabet")

		var alphabet string
		var err error
		if alphabetSet != "custom" {
			alphabet, err = h.uc.GetAlphabet(alphabetSet)
			if err != nil {
				render.Render(w, r, responses.ErrorRenderer(err))
				return
			}
		} else {
			alphabet = alphabet_param
		}

		volume := h.uc.GetInformationVolumeSymbol(alphabet)

		if err := render.Render(w, r, responses.SucceededRenderer(volume)); err != nil {
			h.log.Error("failed succeeded render", logger.Err(err))
			render.Render(w, r, responses.ErrInternalServer)
			return
		}
	}
}

func (h *MeasuringInformationHandler) GetAmountOfInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query().Get("text")
		alphabetSet := r.URL.Query().Get("alphabet_set")
		alphabet_param := r.URL.Query().Get("alphabet")

		var alphabet string
		var err error
		if alphabetSet != "custom" {
			alphabet, err = h.uc.GetAlphabet(alphabetSet)
			if err != nil {
				render.Render(w, r, responses.ErrorRenderer(err))
				return
			}
		} else {
			alphabet = alphabet_param
		}

		amount := h.uc.GetAmountOfInformation(text, alphabet)
		if err := render.Render(w, r, responses.SucceededRenderer(amount)); err != nil {
			h.log.Error("failed succeeded render", logger.Err(err))
			render.Render(w, r, responses.ErrInternalServer)
			return
		}
	}
}
