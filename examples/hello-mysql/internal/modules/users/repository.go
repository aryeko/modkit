package users

import "context"

type Repository interface {
	GetUser(ctx context.Context, id int64) (User, error)
}
