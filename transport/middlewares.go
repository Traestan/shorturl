package transport

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Traestan/shorturl/repository"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/log"
)

func MyMiddleware(logger log.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "OPTIONS" {

			tokenString := r.Header.Get("Authorization")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					logger.Log("msg", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"]))
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return repository.JwtKey, nil
			})

			if err != nil {
				logger.Log("msg", fmt.Errorf("ERROR %v", err))
				return
			}
			if token.Valid {
				logger.Log("msg", "You look nice today")
			} else if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					logger.Log("err", "That's not even a token")
					http.Error(w, "ERROR", http.StatusUnauthorized)
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					logger.Log("err", "Timing is everything")
					http.Error(w, "ERROR", http.StatusUnauthorized)
					return
				} else {
					logger.Log("err", err)
					http.Error(w, "ERROR", http.StatusUnauthorized)
					return
				}
			} else {
				logger.Log("err 3", err)
				http.Error(w, "ERROR", http.StatusUnauthorized)
			}
			ctx := context.WithValue(r.Context(), "Username", token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
