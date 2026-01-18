package main

import (
	"github.com/tuanta7/errx/errors"
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
		if errors.Is(err, errors.ErrRecordNotFound) {
			return nil, errors.ErrRecordNotFound.WithCode(ErrCounterNotFound)
		}
		return nil, err
	}

	return counter, nil
}

func (uc *UseCase) SetCounter(key string, counter *Counter) error {
	err := uc.repo.SetCounter(key, counter)
	if err != nil {
		// adding custom error code
		return err
	}

	return nil
}
