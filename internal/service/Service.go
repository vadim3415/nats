package service

import (
	"nats/internal/model"
	"nats/internal/repository"
)

type PqNats interface {
	PqGetId(id string) (model.ModelNats, error)
	PqNatsMsgCreate(b model.ModelNats) error
}
type Service struct {
	PqNats
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		PqNats: NewUserService(repos.PqNats),
	}
}

///////////////////////////////////////////////////////////////
type UserService struct {
	repo repository.PqNats
}

func NewUserService(repo repository.PqNats) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) PqGetId(id string) (model.ModelNats, error) {
	return s.repo.PqGetId(id)
}

func (s *UserService) PqNatsMsgCreate(b model.ModelNats) error {
	return s.repo.PqNatsMsgCreate(b)
}
