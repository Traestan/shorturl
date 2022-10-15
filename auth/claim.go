package auth

import (
	"context"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
)

const JwtFolder = "jwt/"

//Claims объект данных авторизованного пользователя
type Claims struct {
	UserIdentity UserIdentity `json:"user"`
	rules        map[string]bool
	stdjwt.StandardClaims
}

//UserIdentity объект персональных данных пользователя
type UserIdentity struct {
	UserID   string `json:"user_id"`
	FullName string `json:"fullname"`
}

//HasAccess проверяет имеет ли пользователь правило доступа
func (c *Claims) HasAccess(key string) bool {
	_, ok := c.rules[key]
	return ok
}

// HasAccess более удобный способ проверки правил пользователя
// если хотя бы одно правило присутствует то имеет доступ
func HasAccess(ctx context.Context, rules []string) bool {
	claim, ok := ctx.Value(jwt.JWTClaimsContextKey).(*Claims)
	if !ok {
		return false
	}
	for _, r := range rules {
		if claim.HasAccess(r) {
			return true
		}
	}
	return false
}
