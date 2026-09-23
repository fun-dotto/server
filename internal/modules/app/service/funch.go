package service

import "github.com/fun-dotto/server/internal/modules/app/domain"

type FunchRepository interface {
	GetMenuItems(query domain.MenuItemQuery) ([]domain.MenuItem, error)
}

type FunchService struct {
	repository FunchRepository
}

func NewFunchService(repository FunchRepository) *FunchService {
	return &FunchService{repository: repository}
}

func (s *FunchService) GetMenuItems(query domain.MenuItemQuery) ([]domain.MenuItem, error) {
	return s.repository.GetMenuItems(query)
}
