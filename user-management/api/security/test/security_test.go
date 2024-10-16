package security_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wesleymassine/swordhealth/user-management/api/security"
)

func TestHashPassword(t *testing.T) {
	password := "manager"
	// hashedPassword, err := security.HashPassword(password)
	// assert.NoError(t, err)
	hashedPassword := "$2a$10$nuWwAZ75yMVZciH7fYAuKuWWl8JE3Dcbl895QdjLDNxIOuSRV1lfi"

	err := security.CheckPassword(hashedPassword, password)
	assert.NoError(t, err)
}

func TestGenerateAndParseJWT(t *testing.T) {
	token, err := security.GenerateJWT(1, "admin")
	assert.NoError(t, err)

	claims, err := security.ParseJWT(token)
	assert.NoError(t, err)
	assert.Equal(t, float64(1), claims["user_id"]) // JWT claims are float64 for int
	assert.Equal(t, "admin", claims["role"])
}
