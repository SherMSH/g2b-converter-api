package service

import (
	"converterapi/internal/config"
	"converterapi/internal/repository"
	"converterapi/internal/utils"
)

type G2bServiceIface interface {
	GetReqType() utils.OfflineReqType
	Call() error
}

type Service struct {
	Repository *repository.Repository
}

func New(cfg *config.Configs, repo *repository.Repository) *Service {
	return &Service{
		Repository: repo,
	}
}
