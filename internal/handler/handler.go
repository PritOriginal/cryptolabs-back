package handler

import (
	"log/slog"

	repository "github.com/PritOriginal/cryptolabs-back/internal/repository/alphabet"
	"github.com/PritOriginal/cryptolabs-back/internal/services"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func GetRouter(log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	alphabetRepo := repository.NewAlphabetRepository()
	infoUseCase := services.NewMeasuringInformation(alphabetRepo)
	infoHandler := NewMeasuringInformation(log, *infoUseCase)
	r.Route("/measuring_information", func(r chi.Router) {
		r.Get("/alphabet", infoHandler.GetAlphabet())
		r.Get("/volume", infoHandler.GetInformationVolumeSymbol())
		r.Get("/amount", infoHandler.GetAmountOfInformation())
	})

	compressionService := services.NewCompressionService()
	compressionHandler := NewCompressionHandler(log, compressionService)
	r.Route("/compression", func(r chi.Router) {
		r.Post("/compress", compressionHandler.Compress())
		r.Post("/decompress", compressionHandler.Decompress())
	})

	return r
}
