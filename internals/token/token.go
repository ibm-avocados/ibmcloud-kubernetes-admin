package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var now = time.Now

type Claims struct {
	Provider string `json:"provider"`
	Login    bool   `json:"login"`
	ExtraData  string `json:"extraData"`
}

type _Claims struct {
	Claims
	jwt.StandardClaims
}

func New(claims Claims) (string, error) {
	expirationTime := now().Add(3 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, _Claims{
		claims,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Verify(token string) (*Claims, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		res := Claims{}
		claimJSON, err := json.Marshal(claims)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(claimJSON, &res)
		if err != nil {
			return nil, err
		}
		return &res, nil
	}

	return nil, errors.New("something happened")
}
