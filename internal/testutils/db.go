package testutils

import (
	"finance-crud-app/internal/db"
	"log"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	ConnStr = "postgres://postgres:Password123@localhost:5432/crud_db?sslmode=disable"
	DB      *sqlx.DB
)

type User struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Records   []Record `json:"records"`
}

type Record struct {
	Description string `json:"description"`
	Category    string `json:"category"`
	Amount      int    `json:"amount"`
}

type Data struct {
	TestUser []User `json:"testUser"`
}

func init() {
	// migrate := flag.Bool("migrate", false, "perform database migration")
	// flag.Parse()

	var err error
	DB, err = db.NewPGStorage(ConnStr)
	if err != nil {
		log.Fatalf("Error connecting to TestingDB: %v", err)
	}
	defer DB.Close()
}

func MigrateTestDb() {
	u, _ := url.Parse(ConnStr)
	db := dbmate.New(u)

	err := db.CreateAndMigrate()
	if err != nil {
		panic(err)
	}
}
