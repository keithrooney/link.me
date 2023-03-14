package internal

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func WithAuthentication(next http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		header := req.Header.Get("authorization")
		if header == "" {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("bad request"))
			return
		}
		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("internal server error")
			}
			return secret, nil
		})
		if err == nil && token.Valid {
			next.ServeHTTP(rw, req)
		} else {
			rw.WriteHeader(http.StatusForbidden)
			rw.Write([]byte("forbidden"))
		}
	}
}
