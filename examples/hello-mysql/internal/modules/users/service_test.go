package users

import (
	"context"
	"errors"
	"testing"
)

type fakeRepo struct {
	user User
	err  error
}

func (f fakeRepo) GetUser(ctx context.Context, id int64) (User, error) {
	if f.err != nil {
		return User{}, f.err
	}
	if f.user.ID != id {
		return User{}, errors.New("not found")
	}
	return f.user, nil
}

func TestServiceGetUser_ReturnsUser(t *testing.T) {
	repo := fakeRepo{user: User{ID: 10, Name: "Ada", Email: "ada@example.com"}}
	svc := NewService(repo)

	got, err := svc.GetUser(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.ID != 10 || got.Name != "Ada" || got.Email != "ada@example.com" {
		t.Fatalf("unexpected user: %+v", got)
	}
}

func TestServiceGetUser_PropagatesError(t *testing.T) {
	repo := fakeRepo{err: errors.New("boom")}
	svc := NewService(repo)

	_, err := svc.GetUser(context.Background(), 10)
	if err == nil {
		t.Fatalf("expected error")
	}
}
