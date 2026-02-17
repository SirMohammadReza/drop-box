package encryption

import "golang.org/x/crypto/bcrypt"

type EncryptionService struct{}

func (es *EncryptionService) HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}

	return string(hashPassword), nil
}

func (es *EncryptionService) CheckHash(password, hashValue string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashValue), []byte(password)) == nil
}
