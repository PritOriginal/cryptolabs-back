package handler

import (
	"log/slog"
	"net/http"

	"github.com/PritOriginal/cryptolabs-back/internal/services"
	"github.com/PritOriginal/problem-map-server/pkg/handlers"
	"github.com/PritOriginal/problem-map-server/pkg/responses"
)

type MeasuringInformationHandler struct {
	handlers.BaseHandler
	uc services.MeasuringInformation
}

func NewMeasuringInformation(log *slog.Logger, uc services.MeasuringInformation) *MeasuringInformationHandler {
	return &MeasuringInformationHandler{handlers.BaseHandler{Log: log}, uc}
}

func (h *MeasuringInformationHandler) GetAlphabet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alphabetSet := r.URL.Query().Get("alphabet_set")

		alphabet, err := h.uc.GetAlphabet(alphabetSet, "")
		if err != nil {
			return
		}

		h.Render(w, r, responses.SucceededRenderer(alphabet))
	}
}

func (h *MeasuringInformationHandler) GetInformationVolumeSymbol() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alphabetSet := r.URL.Query().Get("alphabet_set")
		customAlphabet := r.URL.Query().Get("alphabet")

		alphabet, err := h.uc.GetAlphabet(alphabetSet, customAlphabet)
		if err != nil {
			h.RenderError(w, r,
				handlers.HandlerError{Msg: "error get alphabet", Err: err},
				responses.ErrBadRequest,
			)
			return
		}

		volume := h.uc.GetInformationVolumeSymbol(alphabet)

		h.Render(w, r, responses.SucceededRenderer(volume))
	}
}

func (h *MeasuringInformationHandler) GetAmountOfInformation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		text := r.URL.Query().Get("text")
		alphabetSet := r.URL.Query().Get("alphabet_set")
		alphabet_param := r.URL.Query().Get("alphabet")

		alphabet, err := h.uc.GetAlphabet(alphabetSet, alphabet_param)
		if err != nil {
			h.RenderError(w, r,
				handlers.HandlerError{Msg: "error get alphabet", Err: err},
				responses.ErrBadRequest,
			)
			return
		}

		amount := h.uc.GetAmountOfInformation(text, alphabet)
		h.Render(w, r, responses.SucceededRenderer(amount))
	}
}
