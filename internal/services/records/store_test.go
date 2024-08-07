package records_test

import (
	"finance-crud-app/internal/db"
	"finance-crud-app/internal/services/records"
	"finance-crud-app/internal/services/user"
	"finance-crud-app/internal/types"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	recordTestStore *records.RecordsStore
	userTestStore   *user.Store
	testDB          *sqlx.DB
)

func TestMain(m *testing.M) {
	ConnStr := "postgres://postgres:Password123@localhost:5432/crud_db?sslmode=disable"
	testDB, err := db.NewPGStorage(ConnStr)
	if err != nil {
		log.Fatalf("could not connect %v", err)
	}
	defer testDB.Close()
	recordTestStore = records.NewStore(testDB)
	userTestStore = user.NewStore(testDB)

	code := m.Run()
	os.Exit(code)
}

func TestGetUserRecords(t *testing.T) {
	user, err := userTestStore.GetUserByEmail("test1@email.com")
	if err != nil {
		t.Fatalf("could not retrieve test user %v", err)
	}

	userId := strconv.Itoa(user.ID)

	testrecord, err := recordTestStore.CreateUserRecord(userId, types.Record{
		Description: "get user test description",
		Category:    "unique testing category",
		Amount:      3000,
	})

	t.Run("get mutliple user records", func(t *testing.T) {
		userIdInt := strconv.Itoa(user.ID)
		_, err := recordTestStore.GetUserRecords(userIdInt)
		if err != nil {
			t.Errorf("could not retrieve user records %v", err)
		}
	})

	t.Run("get record by id", func(t *testing.T) {
		_, err := recordTestStore.GetRecordById(strconv.Itoa(testrecord))
		if err != nil {
			t.Errorf("could not retrieve user recods %v", err)
		}
	})

	t.Run("get record by category", func(t *testing.T) {

	})

	t.Cleanup(func() {
		record, err := recordTestStore.GetRecordById(strconv.Itoa(testrecord))
		if err != nil {
			t.Fatalf("could not retrieve get record test data %v", err)
		}

		err = recordTestStore.DeleteRecord(strconv.Itoa(record.ID))
		if err != nil {
			t.Errorf("could note delete test record data %v", err)
		}
	})
}

func TestCreateUserRecord(t *testing.T) {
	testUser, err := userTestStore.CreateUser(types.User{
		Password:  "1334543",
		Email:     "testUser@records.com",
		FirstName: "RecordsTestUser",
		LastName:  "RecordsTestUserLastName",
	})
	if err != nil {
		t.Fatalf("test user creation failed with error %v", testUser)
	}

	testUserId := strconv.Itoa(testUser)

	test_data := map[string]struct {
		record_data types.Record
		userId      string
		result      any
	}{
		"PASS create record": {
			record_data: types.Record{
				Description: "example description",
				Amount:      200,
			},
			userId: testUserId,
			result: nil,
		},
	}

	for name, tc := range test_data {
		t.Run(name, func(t *testing.T) {
			_, got := recordTestStore.CreateUserRecord(tc.userId, tc.record_data)
			if got != tc.result {
				t.Errorf("error testing function got %v expected value %v", got, tc.result)
			}
		})
	}

	t.Cleanup(func() {
		err := userTestStore.DeleteUser("testUser@records.com")
		if err != nil {
			t.Errorf("could not delete user %v got error %v", "testUser@records.com", err)
		}
	})
}

func TestUserDeleteRecord(t *testing.T) {
	testUser, err := userTestStore.CreateUser(types.User{
		Password:  "1334543",
		Email:     "testUser@recordstest.com",
		FirstName: "RecordsTestUser",
		LastName:  "RecordsTestUserLastName",
	})
	if err != nil {
		t.Fatalf("test user creation failed with error %v", err)
	}

	testRecord, err := recordTestStore.CreateUserRecord(strconv.Itoa(testUser), types.Record{
		Description: "delete user test record",
		Amount:      2343,
	})
	if err != nil {
		t.Fatalf("test create user record failed with error %v", err)
	}

	t.Run("delete created test user", func(t *testing.T) {
		err := recordTestStore.DeleteRecord(strconv.Itoa(testRecord))
		if err != nil {
			t.Fatalf("test delete user record failed with error %v", err)
		}
	})

}

func TestCheckRecordBelongsToUser(t *testing.T) {
	testUser, err := userTestStore.CreateUser(types.User{
		Password:  "1334543",
		Email:     "testUser@usersrecordstest.com",
		FirstName: "RecordsTestUser",
		LastName:  "RecordsTestUserLastName",
	})
	if err != nil {
		t.Fatalf("test user creation failed with error %v", err)
	}

	testRecord, err := recordTestStore.CreateUserRecord(strconv.Itoa(testUser), types.Record{
		Description: "delete user test record",
		Amount:      2343,
	})
	if err != nil {
		t.Fatalf("test create user record failed with error %v", err)
	}

	t.Run("successfull check record belongs to user", func(t *testing.T) {
		ok := recordTestStore.CheckRecordBelongsToUser(strconv.Itoa(testUser), strconv.Itoa(testRecord))
		if !ok {
			t.Fatalf("test delete user record failed with error %v", err)
		}
	})

	t.Cleanup(func() {
		err := recordTestStore.DeleteRecord(strconv.Itoa(testRecord))
		if err != nil {
			t.Fatalf("could not delete record %v", err)
		}

		err = userTestStore.DeleteUser("testUser@usersrecordstest.com")
		if err != nil {
			t.Fatalf("could not delete record %v", err)
		}
	})
}
