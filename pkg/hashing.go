package pkg

import "golang.org/x/crypto/bcrypt"

func Hash(text string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// hashString := string(hash)
	return string(hash), nil
}
