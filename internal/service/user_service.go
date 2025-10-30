package service

import (
	"context"
	"fmt"
	"go-multitenant/internal/db"
	"go-multitenant/internal/repo"
)

type UserService struct {
	Provider *db.DBProvider
	Repo     *repo.UserRepo
}

func NewUserService(p *db.DBProvider, r *repo.UserRepo) *UserService {
	return &UserService{Provider: p, Repo: r}
}
func (s *UserService) GetUser(ctx context.Context, clientID string, id uint) (interface{}, error) {
	clientDB, err := s.Provider.GetClientDB(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("client db: %w", err)
	}
	u, err := s.Repo.FindByID(clientDB, id)
	if err != nil {
		return nil, err
	}
	return u, nil
}
