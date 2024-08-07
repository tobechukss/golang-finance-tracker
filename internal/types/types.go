package types

type User struct {
	ID        int    `json:"id"`
	Password  string `json:"-"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	CreatedAt string `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) (userId int, err error)
	DeleteUser(email string) error
}

type Record struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Amount      int    `json:"amount"`
	UserId      int    `json:"user_id"`
	CreatedAt   string `json:"createAt"`
}

type RecordStore interface {
	GetUserRecords(userId string) ([]Record, error)
	GetRecordById(id string) (Record, error)
	GetUserRecordsByCategory(userId string, category string) ([]Record, error)
	CreateUserRecord(userId string, record Record) (recordId int, err error)
	CheckRecordBelongsToUser(userId, recordId string) bool
	UserDeleteRecord(recordId, userId string) error
	DeleteRecord(recordId string) error
}

type RegisterUserPayload struct {
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PostRecordPayload struct {
	Description string `json:"description" validate:"required"`
	Category    string `json:"category"`
	Amount      int    `json:"amount" validate:"required"`
}
