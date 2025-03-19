package services

import (
	"reflect"
	"testing"
)

func TestCompressionService_Compress(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want []byte
	}{
		{
			name: "6W",
			data: []byte("WWWWWW"),
			want: []byte("6W"),
		},
		{
			name: "4B",
			data: []byte("BBBB"),
			want: []byte("4B"),
		},
		{
			name: "",
			data: []byte("WWWWBBBWBB"),
			want: []byte("4W3B1W2B"),
		},
		{
			name: "",
			data: []byte("WWBWBWBBWBB"),
			want: []byte("2W1B1W1B1W2B1W2B"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CompressionService{}
			if got := s.Compress(tt.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CompressionService.Compress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompressionService_Decompress(t *testing.T) {

	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "",
			data:    []byte("1W1B"),
			want:    []byte("WB"),
			wantErr: false,
		},
		{
			name:    "",
			data:    []byte("3B2W1B"),
			want:    []byte("BBBWWB"),
			wantErr: false,
		},
		{
			name:    "",
			data:    []byte("2W3B"),
			want:    []byte("WWBBB"),
			wantErr: false,
		},
		{
			name:    "invalid_1",
			data:    []byte("2W3"),
			want:    []byte("WWBBB"),
			wantErr: true,
		},
		{
			name:    "invalid_2",
			data:    []byte("2WB"),
			want:    []byte("WWBBB"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CompressionService{}
			got, err := s.Decompress(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompressionService.Decompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("CompressionService.Decompress() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
