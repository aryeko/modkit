package users

import (
	"context"

	"github.com/aryeko/modkit/examples/hello-mysql/internal/sqlc"
)

type mysqlRepo struct {
	queries *sqlc.Queries
}

func NewMySQLRepo(queries *sqlc.Queries) Repository {
	return &mysqlRepo{queries: queries}
}

func (r *mysqlRepo) GetUser(ctx context.Context, id int64) (User, error) {
	row, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return User{}, err
	}
	return User{ID: row.ID, Name: row.Name, Email: row.Email}, nil
}
