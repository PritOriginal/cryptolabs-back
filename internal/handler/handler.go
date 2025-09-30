package handler

import (
	"log/slog"

	repository "github.com/PritOriginal/cryptolabs-back/internal/repository/alphabet"
	"github.com/PritOriginal/cryptolabs-back/internal/services"
	"github.com/PritOriginal/cryptolabs-back/internal/services/compression"
	"github.com/PritOriginal/cryptolabs-back/internal/services/crypto"
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
	infoService := services.NewMeasuringInformation(alphabetRepo)
	infoHandler := NewMeasuringInformation(log, *infoService)
	r.Route("/measuring_information", func(r chi.Router) {
		r.Get("/alphabet", infoHandler.GetAlphabet())
		r.Get("/volume", infoHandler.GetInformationVolumeSymbol())
		r.Get("/amount", infoHandler.GetAmountOfInformation())
	})

	type CompressionServiceItem struct {
		name    string
		service CompressionService
	}

	compressionServices := []CompressionServiceItem{
		{name: "/rle", service: compression.NewRLEService()},
		{name: "/huffman", service: compression.NewHuffmanService()},
		{name: "/arithmetic", service: compression.NewArithmeticService()},
		{name: "/lzw", service: compression.NewLZWService()},
	}
	for _, serviceItem := range compressionServices {
		serviceHandler := NewCompressionHandler(log, serviceItem.service)
		r.Route(serviceItem.name, func(r chi.Router) {
			r.Post("/compress", serviceHandler.Compress())
			r.Post("/compress/details", serviceHandler.CompressWithDetails())
			r.Post("/decompress", serviceHandler.Decompress())
			r.Post("/decompress/details", serviceHandler.DecompressWithDetails())
		})
	}

	rsaService := crypto.NewRsaService()
	rsaHandler := NewRsaHandler(log, rsaService)
	r.Route("/rsa", func(r chi.Router) {
		r.Get("/keys", rsaHandler.GenerateKeys())
		r.Post("/encrypt", rsaHandler.Encrypt())
		r.Post("/decrypt", rsaHandler.Decrypt())
	})

	return r
}
