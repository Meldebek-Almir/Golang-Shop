package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"shop/backend-api/internal/jwt"
)

type keyUserType string

const keyUser = "user"

func AuthRedirectMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		status, _, err := jwt.VerifyJWT(c.Value)
		if err != nil {
			fmt.Println(err)
		}
		if status {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthNeedMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token")
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		status, claims, err := jwt.VerifyJWT(c.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if !status {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		sub, ok := claims["sub"]
		if !ok {
			log.Println("claims sub not exists")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), keyUserType(keyUser), sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SellerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		status, claims, err := jwt.VerifyJWT(c.Value)
		if !status || err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"]
		if !ok {
			log.Println("claims role not exists")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if role != "seller" && role != "admin" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		sub, ok := claims["sub"]
		if !ok {
			log.Println("claims sub not exists")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), keyUserType(keyUser), sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClientMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		status, claims, err := jwt.VerifyJWT(c.Value)
		if !status || err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if role != "client" && role != "admin" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		sub, ok := claims["sub"]
		if !ok {
			log.Println("claims sub not exists")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), keyUserType(keyUser), sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("jwt_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		status, claims, err := jwt.VerifyJWT(c.Value)
		if !status || err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if role != "admin" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
