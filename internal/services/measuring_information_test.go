package services

import (
	"reflect"
	"testing"

	repository "github.com/PritOriginal/cryptolabs-back/internal/repository/alphabet"
)

func TestMeasuringInformation_GetAlphabet(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &MeasuringInformation{
				alphabetRepo: repository.NewMockAlphabetRepository(t),
			}
			got, err := uc.GetAlphabet(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MeasuringInformation.GetAlphabet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MeasuringInformation.GetAlphabet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeasuringInformation_GetInformationVolumeSymbol(t *testing.T) {

	type args struct {
		alphabet string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test-empty",
			args: args{alphabet: ""},
			want: 0,
		},
		{
			name: "test-1",
			args: args{alphabet: "ab"},
			want: 1,
		},
		{
			name: "test-2",
			args: args{alphabet: "abc"},
			want: 2,
		},
		{
			name: "test-3",
			args: args{alphabet: "абв"},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &MeasuringInformation{
				alphabetRepo: repository.NewMockAlphabetRepository(t),
			}
			if got := uc.GetInformationVolumeSymbol(tt.args.alphabet); got != tt.want {
				t.Errorf("MeasuringInformation.GetInformationVolumeSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMeasuringInformation(t *testing.T) {
	type args struct {
		repo repository.AlphabetRepository
	}
	tests := []struct {
		name string
		args args
		want *MeasuringInformation
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMeasuringInformation(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMeasuringInformation() = %v, want %v", got, tt.want)
			}
		})
	}
}
