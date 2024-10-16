package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = "ace353e1c2dd28f9fa8c40f3687f943f7a4c0576dedc702fc049f7f98f06467a" // Change to a strong key in production

// GenerateJWT generates a new JWT token for a given user ID and role.
func GenerateJWT(userID int64, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

// ParseJWT parses and validates a JWT token, returning the claims if valid.
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return jwtSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}
