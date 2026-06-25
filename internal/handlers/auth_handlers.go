package handlers

import (
	"database/sql" // <-- Pacote nativo adicionado aqui
	"encoding/json"
	"net/http"

	"github.com/boma026/BarberFlow-API/internal/auth"
	"github.com/boma026/BarberFlow-API/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Erro ao processar senha", http.StatusInternalServerError)
			return
		}

		user, err := q.CreateUser(r.Context(), db.CreateUserParams{
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
		})
		if err != nil {
			http.Error(w, "Erro ao criar usuário (email pode já existir)", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Usuário criado com sucesso", "email": user.Email})
	}
}

func LoginUser(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		user, err := q.GetUserByEmail(r.Context(), req.Email)
		if err != nil {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
		if err != nil {
			http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
			return
		}

		accessToken, refreshToken, err := auth.GenerateTokens(user.Email)
		if err != nil {
			http.Error(w, "Erro ao gerar tokens", http.StatusInternalServerError)
			return
		}

		q.UpdateRefreshToken(r.Context(), db.UpdateRefreshTokenParams{
			ID: user.ID,
			// Correção aplicada na linha abaixo:
			RefreshToken: sql.NullString{String: refreshToken, Valid: true},
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func RefreshToken(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}

		claims, err := auth.ValidateToken(req.RefreshToken)
		if err != nil {
			http.Error(w, "Refresh token inválido ou expirado", http.StatusUnauthorized)
			return
		}

		user, err := q.GetUserByEmail(r.Context(), claims.Email)
		if err != nil {
			http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
			return
		}

		if !user.RefreshToken.Valid || user.RefreshToken.String != req.RefreshToken {
			http.Error(w, "Refresh token revogado ou inválido", http.StatusUnauthorized)
			return
		}

		newAccessToken, newRefreshToken, err := auth.GenerateTokens(user.Email)
		if err != nil {
			http.Error(w, "Erro ao rotacionar tokens", http.StatusInternalServerError)
			return
		}

		q.UpdateRefreshToken(r.Context(), db.UpdateRefreshTokenParams{
			ID:           user.ID,
			RefreshToken: sql.NullString{String: newRefreshToken, Valid: true},
		})

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  newAccessToken,
			"refresh_token": newRefreshToken,
		})
	}
}
