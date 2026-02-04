package users

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	modkithttp "github.com/aryeko/modkit/modkit/http"
)

type stubService struct {
	user User
	err  error
}

func (s stubService) GetUser(ctx context.Context, id int64) (User, error) {
	if s.err != nil {
		return User{}, s.err
	}
	return s.user, nil
}

func TestController_GetUser_ReturnsJSON(t *testing.T) {
	ctrl := NewController(stubService{user: User{ID: 7, Name: "Lin", Email: "lin@example.com"}})
	router := modkithttp.NewRouter()
	ctrl.RegisterRoutes(modkithttp.AsRouter(router))

	req := httptest.NewRequest(http.MethodGet, "/users/7", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var got User
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if got.ID != 7 || got.Name != "Lin" || got.Email != "lin@example.com" {
		t.Fatalf("unexpected user: %+v", got)
	}
}
