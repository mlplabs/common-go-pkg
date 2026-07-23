package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

var TempSecret = "123"

// https://levelup.gitconnected.com/crud-restful-api-with-go-gorm-jwt-postgres-mysql-and-testing-460a85ab7121

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ValidPassword(hashedPassword []byte, password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		log.Printf("Error compare hash and pssword. %v", err)
		err = errors.New("wrong username or password")
	}

	return err == nil, err
}

func CreateTokenPair(claims map[string]any, secretKey string,
	tokenExpiresSec int64, refreshExpiresSec int64,
) (*TokenPair, error) {

	accessClaims := jwt.MapClaims{}
	accessClaims["exp"] = time.Now().Add(time.Duration(tokenExpiresSec) * time.Second).Unix()
	refreshClaims := jwt.MapClaims{}
	refreshClaims["exp"] = time.Now().Add(time.Duration(refreshExpiresSec) * time.Second).Unix()

	for k, v := range claims {
		accessClaims[k] = v
		refreshClaims[k] = v
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Create signed access tocken string %v", err)
	}
	refreshTokenSting, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		log.Printf("Create signed refresh tocken string %v", err)
	}

	return &TokenPair{accessTokenString, "bearer", tokenExpiresSec, refreshTokenSting}, err
}

func ReadToken(secretKey string, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return token, nil
	}

	return nil, err
}

func ReadTokenUnverified(tokenString string) (*jwt.Token, []string, error) {
	token, parts, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, nil, err
	}

	return token, parts, nil
}

func TokenValid(r *http.Request, secretKey string) error {
	tokenString := ExtractToken(r)

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}

	return nil
}

func ExtractToken(r *http.Request) string {
	keys := r.URL.Query()
	token := keys.Get("token")
	if token != "" {
		return token
	}
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}

	return ""
}
