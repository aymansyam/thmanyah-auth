package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type loginReq struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func main() {
	port := env("PORT", "8081")
	secret := env("JWT_SECRET", "dev-secret")

	mux := http.NewServeMux()

	// Liveness
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Issue a JWT for any non-empty user/password (demo only)
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var body loginReq
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.User == "" || body.Password == "" {
			http.Error(w, "bad credentials", http.StatusUnauthorized)
			return
		}
		claims := jwt.MapClaims{
			"sub": body.User,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(1 * time.Hour).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(secret))
		if err != nil {
			http.Error(w, "signing error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": signed})
	})

	// Validate a JWT sent in Authorization: Bearer <token>
	mux.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		raw := auth
		if len(auth) > 7 && auth[:7] == "Bearer " {
			raw = auth[7:]
		}
		_, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("valid"))
	})

	log.Printf("auth listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
