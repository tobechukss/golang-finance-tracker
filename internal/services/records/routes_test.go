package records

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"finance-crud-app/internal/types"

	"github.com/gorilla/mux"
)

func TestRecordsServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	recordStore := &mockRecordStore{}

	userID := "12345"

	handler := NewHandler(recordStore, userStore)

	t.Run("should pass retrieve a list of records", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/record", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/record", handler.handleGetRecord).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusAccepted {
			t.Errorf("expected status code %d, got %d", http.StatusAccepted, rr.Code)
		}
	})

	t.Run("should pass creating a record with correct record payload", func(t *testing.T) {
		record := map[string]any{
			"description": "value",
			"amount":      100,
		}
		jsonPayload, err := json.Marshal(record)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(context.Background(), "user_id", userID)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, "/record", bytes.NewBuffer(jsonPayload))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/record", handler.handlePostRecord).Methods(http.MethodPost)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusAccepted, rr.Code)
		}
	})

	t.Run("should pass retrieving record with recordId", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "user_id", userID)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/record/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/record/{recordID}", handler.handleGetRecordById).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusAccepted {
			t.Errorf("expected status code %d, got %d", http.StatusAccepted, rr.Code)
		}
	})

	t.Run("should pass delete a record with correct record payload", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/record/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/record/{recordID}", handler.handleDeleteRecord).Methods(http.MethodGet)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusAccepted {
			t.Errorf("expected status code %d, got %d", http.StatusAccepted, rr.Code)
		}
	})
}

type mockRecordStore struct{}

func (m *mockRecordStore) GetUserRecords(userId string) ([]types.Record, error) {
	return []types.Record{}, nil
}

func (m *mockRecordStore) GetRecordById(id string) (types.Record, error) {
	return types.Record{}, nil
}

func (m *mockRecordStore) GetUserRecordsByCategory(userId string, category string) ([]types.Record, error) {
	return []types.Record{}, nil
}

func (m *mockRecordStore) CreateUserRecord(userId string, record types.Record) (recordId int, err error) {
	return 1, nil
}

func (m *mockRecordStore) UserDeleteRecord(recordId, userId string) error {
	return nil
}

func (m *mockRecordStore) DeleteRecord(recordId string) error {
	return nil
}

func (m *mockRecordStore) CheckRecordBelongsToUser(userId, recordId string) bool {
	return true
}

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
