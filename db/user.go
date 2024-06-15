package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"-"`
	IsAdmin   bool       `gorm:"default:false" json:"isAdmin"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// CreateAdmin is a method on the User struct that creates an admin user in the database.
// It creates a new User object with the specified email and password and sets the IsAdmin field to true.
// It then hashes the password and saves the user to the database.
//
// This method does not take any parameters.
//
// Returns:
// error: An error object that describes an error that occurred during the method's execution.
func (u *User) CreateAdmin() error {
	user := User{
		Email:    "your email",
		Password: "your password",
		IsAdmin:  true,
	}
	// Hash Password & update user
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return errors.New("error creating password")
	}
	user.Password = string(password)
	// Create User
	if err := DBConn.Create(&user).Error; err != nil {
		return errors.New("error creating user")
	}
	return nil
}

// LoginAsAdmin is a method on the User struct that logs in a user as an admin.
// It takes an email and password as input and attempts to find a user with the specified email and is_admin set to true.
// If the user is found, it compares the password hash with the input password and returns the user if the passwords match.
// If the user is not found or the passwords do not match, it returns an error.
//
// Parameters:
// email string: The email of the user to log in.
// password string: The password of the user to log in.
//
// Returns:
// *User: A pointer to the User object representing the logged-in user.
// error: An error object that describes an error that occurred during the method's execution.
func (u *User) LoginAsAdmin(email string, password string) (*User, error) {
	// Find User
	if err := DBConn.Where("email = ? AND is_admin = ?", email, true).First(&u).Error; err != nil {
		return nil, errors.New("user not found")
	}
	// Compare Passwords
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return u, nil
}
