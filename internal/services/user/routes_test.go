package user

import (
	"finance-crud-app/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type mockUserStore struct{}

func (m *mockUserStore) DeleteUser(email string) error {
	return nil
}

func (m *mockUserStore) GetUserByEmail(email string) (types.User, error) {
	return types.User{}, nil
}

func (m *mockUserStore) CreateUser(u types.User) (userId int, err error) {
	return 1, nil
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return &types.User{}, nil
}

func TestGetUserHandler(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail to get user with user_id that is not a number", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/user/abc", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}", handler.handleGetUser).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should pass to get user with numeric id", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/user/23", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/user/{userID}", handler.handleGetUser).Methods(http.MethodGet)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusBadGateway, rr.Code)
		}
	})
}
