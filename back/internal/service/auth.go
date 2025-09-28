package service

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"challenge_enube/internal/models"
	"challenge_enube/internal/usecase"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailInUse         = errors.New("email already in use")
	ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
	userRepo       *models.UserRepository
	jwtSecret      []byte
	accessTokenTTL time.Duration
}

func NewAuthService(userRepo *models.UserRepository, jwtSecret string, accessTokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		jwtSecret:      []byte(jwtSecret),
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *AuthService) Register(email, username, password string) (*models.User, error) {
	_, err := s.userRepo.GetUserByEmail(email)
	if err == nil {
		return nil, ErrEmailInUse
	}

	hashedPassword, err := usecase.HashPassword(password)
	if err != nil {
		return nil, err
	}

	return s.userRepo.CreateUser(email, username, hashedPassword)
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := usecase.VerifyPassword(user.PasswordHash, password); err != nil {
		return "", ErrInvalidCredentials
	}

	expiration := time.Now().Add(s.accessTokenTTL)
	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.Username,
		"email":    user.Email,
		"exp":      expiration.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}

		_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return s.jwtSecret, nil
		})

		if err != nil {
			http.Error(w, "token inválido ou expirado", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
