package jwt

import (
	"fmt"

	jwtp "github.com/golang-jwt/jwt/v4"
	"github.com/melvinodsa/go-iam/sdk"
)

type service struct {
	secret []byte
}

func NewService(secret sdk.MaskedBytes) Service {
	return &service{secret: []byte(secret)}
}

// GenerateToken generates a JWT token with the given claims
func (s service) GenerateToken(claims map[string]interface{}, expiryTimeInSeconds int64) (string, error) {
	jwtClaims := jwtp.MapClaims{}
	for k, v := range claims {
		jwtClaims[k] = v
	}
	jwtClaims["exp"] = expiryTimeInSeconds
	token := jwtp.NewWithClaims(jwtp.SigningMethodHS256, jwtClaims)

	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("error signing jwt token with secret %w", err)
	}
	return tokenString, nil
}

// ValidateToken validates the given JWT token and returns the claims
func (s service) ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwtp.Parse(tokenString, func(token *jwtp.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwtp.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing jwt token %w", err)
	}
	result := map[string]interface{}{}
	if claims, ok := token.Claims.(jwtp.MapClaims); ok && token.Valid {
		for k, v := range claims {
			if k == "exp" {
				continue
			}
			result[k] = v
		}
		return result, nil
	}
	return result, nil
}
