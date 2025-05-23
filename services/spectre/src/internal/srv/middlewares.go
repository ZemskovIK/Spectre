package server

import (
	"context"
	"net/http"
	"spectre/internal/lib"
	"spectre/internal/srv/auth"
	"spectre/pkg/logger"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ! TODO : error handling

func JSONRespMiddleware(nx http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		nx.ServeHTTP(w, r)
	})
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.URL.Path == "/login" {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}
// 		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return []byte("test_secret"), nil // ! TODO
// 		})
// 		if err != nil || !token.Valid {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 			if expRaw, ok := claims["exp"]; ok {
// 				exp, ok := expRaw.(float64)
// 				if ok {
// 					expTime := time.Unix(int64(exp), 0)
// 					if time.Now().After(expTime) {
// 						http.Error(w, "Token expired", http.StatusUnauthorized)
// 						return
// 					}
// 				}
// 			}
// 			ctx := context.WithValue(r.Context(), lib.UserIDKey, claims["sub"])
// 			ctx = context.WithValue(ctx, lib.UserAccessLevelKey, claims["role"])
// 			r = r.WithContext(ctx)
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

func AuthMiddleware(logger *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == auth.LOGIN_POINT {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte("test_secret"), nil // ! TODO
			})
			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if expRaw, ok := claims["exp"]; ok {
					exp, ok := expRaw.(float64)
					if ok {
						expTime := time.Unix(int64(exp), 0)
						if time.Now().After(expTime) {
							http.Error(w, "Token expired", http.StatusUnauthorized)
							return
						}
					}
				}
				ctx := context.WithValue(r.Context(), lib.UserIDKey, claims["sub"])
				ctx = context.WithValue(ctx, lib.UserAccessLevelKey, claims["role"])
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}
