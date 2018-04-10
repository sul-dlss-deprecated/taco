package authorization

import (
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/errors"
)

// AuthenticationService is a function definition that recieves a header
// and returns an Agent or an error
type AuthenticationService func(string) (*Agent, error)

// Authentication returns an authentication service that uses JWT
func Authentication(secretKey string) AuthenticationService {
	return func(header string) (*Agent, error) {
		parts := strings.Split(header, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return nil, errors.New(401, "Invalid Authorization header: %s", header)
		}

		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New(401, "Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			sub, ok := claims["sub"]
			if !ok {
				return nil, errors.New(401, "Sub not provided")
			}
			iss, ok := claims["iss"]
			if !ok {
				return nil, errors.New(401, "Iss not provided")
			}
			return &Agent{Identifier: sub.(string), Service: iss.(string)}, nil
		}

		if err.(*jwt.ValidationError).Errors&jwt.ValidationErrorExpired != 0 {
			return nil, errors.New(401, "Token is expired")
		}
		return nil, err
	}
}
