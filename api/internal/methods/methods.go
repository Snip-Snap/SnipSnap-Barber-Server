package methods

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// CheckError checks if err returns nil. If so, prints and calls os.Exit(1).
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// HashPassword hashes a passed rawPassword.
func HashPassword(rawPassword string) (string, error) {
	cost := bcrypt.DefaultCost
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), cost)
	return string(bytes), err
}

//CheckPasswordHash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// WrongUsernameOrPasswordError is used to deal with authentication erros.
type WrongUsernameOrPasswordError struct{}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}