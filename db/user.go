package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string     `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email    string     `gorm:"unique" json:"email"`
	Password string     `json:"-"` // Don't return password in JSON
	IsAdmin  bool       `gorm:"default:false" json:"isAdmin"`
	CreateAt *time.Time `gorm:"autoCreateTime" json:"createAt"`
	UpdateAt time.Time  `gorm:"autoUpdateTime" json:"updateAt"`
}

func (u *User) CreateAdmin() error {
	user := User{
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  true,
	}

	// Hash password
	// salt: number of rounds to use when hashing the password
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)

	if err != nil {
		return errors.New("error creating password") // error string should not be capitalized
	}

	user.Password = string(password)
	// Create user in db

	if err := DBConn.Create(&user).Error; err != nil {
		return errors.New("error creating user")
	}

	return nil
}

func (u *User) LoginAsAdmin(email string, password string) (*User, error) {
	//Find
	if err := DBConn.Where("email = ? AND is_admin = ?", email, true).First(&u).Error; err != nil {
		// nil: empty value
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(u.Password), /*password in DB*/
		[]byte(password) /*input password*/); err != nil {
		return nil, errors.New("invalid password")
	}
	return u, nil

}
