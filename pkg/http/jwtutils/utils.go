package jwtutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// JSON is a map alias.
type JSON map[string]interface{}

// RenderJSONWithStatus sends data as json and enforces status code.
func RenderJSONWithStatus(w http.ResponseWriter, data interface{}, code int) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}

var (
	ErrKeyMustBePEMEncoded     = errors.New("Invalid Key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	ErrNotRSAPrivateKey        = errors.New("Key is not a valid RSA private key")
	ErrNotRSAPublicKey         = errors.New("Key is not a valid RSA public key")
	ErrUnexpectedSigningMethod = errors.New("Unexpected signing method")
)

type handler struct {
	handler http.Handler
	key     string
}

// TokenValidate is a middleware for jwt token validation.
func TokenValidate(key string) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		handler := handler{key: key, handler: h}
		return validate(handler)
	}
	return f
}

func validate(h handler) http.Handler {
	return http.HandlerFunc(h.middlewareFuncJWT)
}

func (h *handler) middlewareFuncJWT(w http.ResponseWriter, r *http.Request) {
	tokenString := ExtractToken(r)
	if tokenString == "" {
		RenderJSONWithStatus(w, JSON{"error": "требуется авторизация"}, http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token validator: %w, alg: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return []byte(h.key), nil
	})
	if err != nil {
		RenderJSONWithStatus(w, JSON{"error": err.Error()}, http.StatusUnauthorized)
		return
	}
	if !token.Valid {
		RenderJSONWithStatus(w, JSON{"error": "невалидный токен"}, http.StatusUnauthorized)
		return
	}
	ctx := r.Context()
	h.handler.ServeHTTP(w, r.WithContext(ctx))
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func CreateTokenPair(payload map[string]any, secretKey string, accExpSec int64, refExpSec int64) (*TokenPair, error) {
	accessClaims := jwt.MapClaims{}
	accessClaims["authorized"] = true
	accessClaims["exp"] = time.Now().Add(time.Duration(accExpSec) * time.Second).Unix() //Token expires after 15 minutes

	refreshClaims := jwt.MapClaims{}
	refreshClaims["exp"] = time.Now().Add(time.Duration(refExpSec) * time.Second).Unix() //Token expires after 12 hour

	if payload != nil {
		for k, v := range payload {
			accessClaims[k] = v
			refreshClaims[k] = v
		}
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("create signed access token string %v", err)
	}
	refreshTokenSting, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("create signed refresh token string %v", err)
	}
	return &TokenPair{accessTokenString, "bearer", accExpSec, refreshTokenSting}, nil
}
