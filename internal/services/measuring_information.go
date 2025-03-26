package services

import (
	"math"
	"unicode/utf8"

	repository "github.com/PritOriginal/cryptolabs-back/internal/repository/alphabet"
)

type MeasuringInformationService interface {
	GetAlphabet(setName string, customAlphabet string) (string, error)
	GetAmountOfInformation(text string, alphabet string) int
	GetInformationVolumeSymbol(alphabet string) int
}

type MeasuringInformation struct {
	alphabetRepo repository.AlphabetRepository
}

func NewMeasuringInformation(repo repository.AlphabetRepository) *MeasuringInformation {
	return &MeasuringInformation{alphabetRepo: repo}
}

func (uc *MeasuringInformation) GetAlphabet(setName string, customAlphabet string) (string, error) {
	var alphabet string
	if setName != "custom" {
		var err error
		alphabet, err = uc.alphabetRepo.Get(setName)
		if err != nil {
			return alphabet, err
		}
	} else {
		alphabet = customAlphabet
	}
	return alphabet, nil
}
func (uc *MeasuringInformation) GetInformationVolumeSymbol(alphabet string) int {
	var volume int
	lenAlphabet := utf8.RuneCountInString(alphabet)
	if lenAlphabet > 0 {
		volume = int(math.Ceil(math.Log2(float64(lenAlphabet))))
	} else {
		volume = 0
	}
	return volume
}

func (uc *MeasuringInformation) GetAmountOfInformation(text string, alphabet string) int {
	alphabet_map := make(map[rune]int)
	for _, ch := range alphabet {
		alphabet_map[ch] = 0
	}

	var amount int
	for _, ch := range text {
		if _, ok := alphabet_map[ch]; ok {
			alphabet_map[ch]++
		}
	}

	count_ch := 0
	for _, count := range alphabet_map {
		count_ch += count
	}

	power := uc.GetInformationVolumeSymbol(alphabet)
	amount = count_ch * power
	return amount
}
