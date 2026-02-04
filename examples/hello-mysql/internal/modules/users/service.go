package users

import "context"

type Service interface {
	GetUser(ctx context.Context, id int64) (User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetUser(ctx context.Context, id int64) (User, error) {
	return s.repo.GetUser(ctx, id)
}
