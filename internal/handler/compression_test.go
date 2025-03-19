package handler

import (
	"log/slog"
	"net/http"
	"reflect"
	"testing"

	"github.com/PritOriginal/cryptolabs-back/internal/services"
)

func TestCompressionHandler_Compress(t *testing.T) {
	type fields struct {
		log *slog.Logger
		s   services.Compression
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &CompressionHandler{
				log: tt.fields.log,
				s:   tt.fields.s,
			}
			if got := h.Compress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompressionHandler.Compress() = %v, want %v", got, tt.want)
			}
		})
	}
}
