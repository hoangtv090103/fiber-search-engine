package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthClaim struct {
	Id    string `json:"id"`
	User  string `json:"user"`
	Admin bool   `json:"role"`
	jwt.RegisteredClaims
}

// CreateNewAuthToken is a function that creates a new JWT (JSON Web Token) for authentication.
// It takes a user ID, email, and admin status as parameters and creates an AuthClaim with these values and some registered claims.
// The registered claims include an expiry date set to 24 hours from now and an issuer.
// The function then creates a new token with these claims and signs it with a secret key from the environment variables.
// If the secret key is not found in the environment variables, the function panics.
// If there is an error signing the token, the function returns an empty string and an error.
// If the token is successfully signed, the function returns the signed token and nil for the error.
//
// Parameters:
// id string: The user ID to include in the AuthClaim.
// email string: The user email to include in the AuthClaim.
// isAdmin bool: The user admin status to include in the AuthClaim.
//
// Returns:
// string, error: The signed token and an error object that describes an error that occurred during the function's execution.
func CreateNewAuthToken(id string, email string, isAdmin bool) (string, error) {
	claims := AuthClaim{
		Id:    id,
		User:  email,
		Admin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "searchengine.com",
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey, exists := os.LookupEnv(("SECRET_KEY"))

	if !exists {
		panic("SECRET_KEY not found in .env")
	}

	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", errors.New("Failed to sign token")
	}

	return signedToken, nil // nil: no error

}
