package service

import "converterapi/internal/repository"

type Service struct {
	Repository *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Repository: repo,
	}
}
