package main

import (
	"errors"

	"github.com/tuanta7/errx"
)

type UseCase struct {
	repo *Repository
}

func NewUseCase(repo *Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetCounter(key string) (*Counter, error) {
	counter, err := uc.repo.GetCounter(key)
	if err != nil {
		if errors.Is(err, errx.ErrRecordNotFound) {
			return nil, errx.ErrRecordNotFound.WithCode(ErrCounterNotFound)
		}
		return nil, err
	}

	return counter, nil
}
