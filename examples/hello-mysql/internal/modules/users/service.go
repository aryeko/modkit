package users

import "context"

type Service interface {
	GetUser(ctx context.Context, id int64) (User, error)
	CreateUser(ctx context.Context, input CreateUserInput) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	UpdateUser(ctx context.Context, id int64, input UpdateUserInput) (User, error)
	DeleteUser(ctx context.Context, id int64) error
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

func (s *service) CreateUser(ctx context.Context, input CreateUserInput) (User, error) {
	return s.repo.CreateUser(ctx, input)
}

func (s *service) ListUsers(ctx context.Context) ([]User, error) {
	return s.repo.ListUsers(ctx)
}

func (s *service) UpdateUser(ctx context.Context, id int64, input UpdateUserInput) (User, error) {
	return s.repo.UpdateUser(ctx, id, input)
}

func (s *service) DeleteUser(ctx context.Context, id int64) error {
	return s.repo.DeleteUser(ctx, id)
}
