package db

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func NewPGStorage(datasource string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", datasource)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type user struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Email     string   `json:"email"`
	Password  string   `json:"password"`
	Records   []record `json:"records"`
}

type record struct {
	Description string `json:"description"`
	Category    string `json:"category"`
	Amount      int    `json:"amount"`
}

type Data struct {
	TestUser []user `json:"testUser"`
}

func SeedTestDB(db *sqlx.DB) {
	seed_data, err := readJsonFile("test_data/test_seed_data.json")
	if err != nil {
		log.Fatalf("error retrieving data from file %v", seed_data)
	}

	for _, user := range seed_data.TestUser {
		create_user_query := `
		INSERT INTO users
		(firstName, lastName, email, password)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

		create_record_query := `
		INSERT INTO records
		(description, category, amount, userid)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

		var userId int
		err = db.QueryRow(create_user_query, user.FirstName, user.LastName, user.Email, user.Password).Scan(&userId)
		if err != nil {
			log.Printf("error seeding user %v", err)
			return
		}

		for _, record := range user.Records {
			_, err = db.Exec(create_record_query, record.Description, record.Category, record.Amount, userId)
			if err != nil {
				log.Printf("error seeding record %v", err)
			}
		}
	}
}

func readJsonFile(filename string) (*Data, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("error with opening file %v", err)
		return nil, err
	}
	defer file.Close()

	file_bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("error with reading bytes %v", err)
		return nil, err
	}

	var data Data

	err = json.Unmarshal(file_bytes, &data)
	if err != nil {
		log.Printf("error with marshalling %v", err)
		return nil, err
	}

	return &data, nil
}
