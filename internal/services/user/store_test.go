package user_test

import (
	"finance-crud-app/internal/db"
	"finance-crud-app/internal/services/user"
	"finance-crud-app/internal/types"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	userTestStore *user.Store
	testDB        *sqlx.DB
)

func TestMain(m *testing.M) {
	// database
	ConnStr := "postgres://postgres:Password123@localhost:5432/crud_db?sslmode=disable"
	testDB, err := db.NewPGStorage(ConnStr)
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer testDB.Close()
	userTestStore = user.NewStore(testDB)

	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	test_data := map[string]struct {
		user   types.User
		result any
	}{
		"should PASS valid user email used": {
			user: types.User{
				FirstName: "testfirsjjlkjt-1",
				LastName:  "testlastkjh-1",
				Email:     "validuser@email.com",
				Password:  "00000000",
			},
			result: nil,
		},
		"should FAIL invalid user email used": {
			user: types.User{
				FirstName: "testFirstName1",
				LastName:  "testLastName1",
				Email:     "test1@email.com",
				Password:  "800890",
			},
			result: user.CreateUserError,
		},
	}

	for name, tc := range test_data {
		t.Run(name, func(t *testing.T) {
			value, got := userTestStore.CreateUser(tc.user)
			if got != tc.result {
				t.Errorf("test fail expected %v got %v instead and value %v", tc.result, got, value)
			}
		})
	}

	t.Cleanup(func() {
		err := userTestStore.DeleteUser("validuser@email.com")
		if err != nil {
			t.Errorf("could not delete user %v got error %v", "validuser@email.com", err)
		}
	})
}

func TestGetUserByEmail(t *testing.T) {
	test_data := map[string]struct {
		email  string
		result any
	}{
		"should pass valid user email address used": {
			email:  "test1@email.com",
			result: nil,
		},
		"should fail invalid user email address used": {
			email:  "validuser@email.com",
			result: user.RetrieveUserError,
		},
	}

	for name, tc := range test_data {
		got, err := userTestStore.GetUserByEmail(tc.email)
		if err != tc.result {
			t.Errorf("test fail expected %v instead got %v", name, got)
		}
	}
}

func TestGetUserById(t *testing.T) {
	testUserId, err := userTestStore.CreateUser(types.User{
		FirstName: "userbyid",
		LastName:  "userbylast",
		Email:     "unique_email",
		Password:  "unique_password",
	})
	if err != nil {
		log.Panicf("got %v when creating testuser", testUserId)
	}

	test_data := map[string]struct {
		user_id int
		result  any
	}{
		"should pass valid user id used": {
			user_id: testUserId,
			result:  nil,
		},
		"should fail invalid user id used": {
			user_id: 0,
			result:  user.RetrieveUserError,
		},
	}

	for name, tc := range test_data {
		t.Run(name, func(t *testing.T) {
			_, got := userTestStore.GetUserByID(tc.user_id)
			if got != tc.result {
				t.Errorf("error retrieving user by id got %v want %v", got, tc.result)
			}
		})
	}

	t.Cleanup(func() {
		err := userTestStore.DeleteUser("unique_email")
		if err != nil {
			t.Errorf("could not delete user %v got error %v", "unique_email", err)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	testUserId, err := userTestStore.CreateUser(types.User{
		FirstName: "userbyid",
		LastName:  "userbylast",
		Email:     "delete_user@email.com",
		Password:  "unique_password",
	})
	if err != nil {
		log.Panicf("got %v when creating testuser", testUserId)
	}

	test_data := map[string]struct {
		user_email string
		result     error
	}{
		"should pass user email address used": {
			user_email: "delete_user@email.com",
			result:     nil,
		},
	}

	for name, tc := range test_data {
		t.Run(name, func(t *testing.T) {
			err = userTestStore.DeleteUser(tc.user_email)
			if err != tc.result {
				t.Errorf("error deletig user got %v instead of %v", err, tc.result)
			}
		})
	}

	t.Cleanup(func() {
		err := userTestStore.DeleteUser("delete_user@email.com")
		if err != nil {
			log.Printf("could not delete user %v got error %v", "delete_user@email.com", err)
		}
	})
}
