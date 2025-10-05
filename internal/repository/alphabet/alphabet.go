package repository

import (
	"fmt"
	"os"
)

type AlphabetRepository interface {
	Get(name string) (string, error)
}

type AlphabetRepo struct {
}

func NewAlphabetRepository() *AlphabetRepo {
	return &AlphabetRepo{}
}

func (repo *AlphabetRepo) Get(name string) (string, error) {
	var alphabet string

	fContent, err := os.ReadFile(fmt.Sprintf("internal/repository/alphabet/data/%s.txt", name))
	if err != nil {
		return alphabet, err
	}
	alphabet = string(fContent)

	return alphabet, nil
}
