package helper

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt"
)

var SECRETKEY []byte = []byte(os.Getenv("SECRET"))

func CreateToken(email string, role string) (string, error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email, // Subject (user identifier)
		"role":  role,
		"exp":   time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat":   time.Now().Unix(),                // Issued at
	})

	tokenString, err := claims.SignedString(SECRETKEY)
	if err != nil {
        log.Error("Error creating token: ", err)
		return "", err
	}

	// Print information about the created token
	log.Info("Token claims added: %+v\n", claims)
	return tokenString, nil
}

// Function to verify JWT tokens
func VerifyToken(tokenString string) (interface{}, error){
	// Parse the token with the secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SECRETKEY, nil
	})

	// Check for verification errors
	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

    //extract the claims
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["email"].(string)
	role := claims["role"].(string)

	// Return the verified token
    return map[string]string{"email": userID, "role": role}, nil
}
