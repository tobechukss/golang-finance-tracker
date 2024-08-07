package records

import (
	"finance-crud-app/internal/types"
	"fmt"
	"strconv"

	"log"

	"github.com/jmoiron/sqlx"
)

type RecordsStore struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *RecordsStore {
	return &RecordsStore{db: db}
}

func (s *RecordsStore) GetUserRecords(userId string) ([]types.Record, error) {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, fmt.Errorf("error converting id to int: %v", err)
	}

	records := []types.Record{}

	err = s.db.Select(&records, "SELECT * FROM records WHERE userId = $1", userIdInt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving records: %v", err)
	}

	return records, nil
}

func (s *RecordsStore) GetRecordById(id string) (types.Record, error) {
	recordId, err := strconv.Atoi(id)
	if err != nil {
		return types.Record{}, fmt.Errorf("error converting id to int: %v", err)
	}

	record := types.Record{}
	err = s.db.Get(&record, "SELECT * FROM records where id = $1 ", recordId)
	if err != nil {
		return types.Record{}, fmt.Errorf("error retrieving record: %v", err)
	}

	return record, nil
}

func (s *RecordsStore) GetUserRecordsByCategory(userId string, category string) ([]types.Record, error) {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return nil, fmt.Errorf("error converting id to int: %v", err)
	}

	records := []types.Record{}

	query := `SELECT * FROM records WHERE userId = ? AND category = ? ORDER BY createdAt`
	err = s.db.Select(&records, query, userIdInt, category)
	if err != nil {
		return nil, fmt.Errorf("error retrieving records: %v", err)
	}

	return records, nil
}

func (s *RecordsStore) CreateUserRecord(userId string, record types.Record) (int, error) {
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return -1, err
	}

	query := `
	INSERT INTO records
	(description, category, amount, userId)
	VALUES ($1, $2, $3, $4)
	returning id`

	var recordId int
	err = s.db.QueryRow(query,
		record.Description, record.Category, record.Amount, userIdInt).Scan(&recordId)
	if err != nil {
		return -1, err
	}

	return recordId, nil
}

func (s *RecordsStore) CheckRecordBelongsToUser(userId string, recordId string) bool {
	query := `
	SELECT * FROM records
	WHERE
		userId = $1 AND id = $2
	LIMIT
		1`

	var record types.Record
	err := s.db.Get(&record, query, userId, recordId)
	if err != nil {
		log.Printf("error value %v", err)
		return false
	}

	return true
}

func (s *RecordsStore) UserDeleteRecord(recordId, userId string) error {
	ok := s.CheckRecordBelongsToUser(userId, recordId)
	if !ok {
		return fmt.Errorf("user cannot delete record")
	}

	err := s.DeleteRecord(recordId)
	if err != nil {
		return err
	}

	return nil
}

func (s *RecordsStore) DeleteRecord(recordId string) error {
	query := `DELETE FROM records WHERE id = $1`

	_, err := s.db.Exec(query, recordId)
	if err != nil {
		return err
	}

	return nil
}
