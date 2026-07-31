package jwtutils

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	domainErr "github.com/mlplabs/common-go-pkg/pkg/http/errors"
	"github.com/mlplabs/common-go-pkg/pkg/http/errors/custom"

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

type CtxId struct{}
type CtxUserId struct{}
type CtxOwnerId struct{}

// JSON is a map alias.
type JSON map[string]interface{}

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
		domainErr.SetError(w, r, custom.NewUnauthorized(errors.New("no token found")))
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("token validator: %w, alg: %v", ErrUnexpectedSigningMethod, token.Header["alg"])
		}
		return []byte(h.key), nil
	})
	if err != nil {
		domainErr.SetError(w, r, custom.NewUnauthorized(err))
		return
	}
	if !token.Valid {
		domainErr.SetError(w, r, custom.NewUnauthorized(errors.New("invalid token")))
		return
	}
	ctx := r.Context()
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id, ok := claims["id"]
		if ok {
			ctx = context.WithValue(ctx, CtxId{}, id)
		}

		userId, ok := claims["user_id"]
		if ok {
			ctx = context.WithValue(ctx, CtxUserId{}, userId)
		}

		ownerId, ok := claims["owner_id"]
		if ok {
			ctx = context.WithValue(ctx, CtxOwnerId{}, ownerId)
		}
	}
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

func GetContextId(r *http.Request) (int, error) {
	uid, ok := r.Context().Value(CtxId{}).(int)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}

	return uid, nil
}

func GetUserId(r *http.Request) (int, error) {
	uid, ok := r.Context().Value(CtxUserId{}).(int)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}

	return uid, nil
}

func GetOwnerId(r *http.Request) (int, error) {
	uid, ok := r.Context().Value(CtxOwnerId{}).(int)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}

	return uid, nil
}
