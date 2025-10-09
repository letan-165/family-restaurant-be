package models

type UserRole string

const (
	ADMIN    UserRole = "ADMIN"
	CUSTOMER UserRole = "CUSTOMER"
)

type User struct {
	ID      string   `bson:"_id,omitempty" json:"id"`
	Name    string   `bson:"name" json:"name"`
	Email   string   `bson:"email" json:"email"`
	Picture string   `bson:"picture" json:"picture"`
	Role    UserRole `bson:"role" json:"role"`
}

func (u UserRole) IsValid() bool {
	switch u {
	case ADMIN, CUSTOMER:
		return true
	default:
		return false
	}
}
