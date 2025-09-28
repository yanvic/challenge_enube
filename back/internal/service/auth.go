package service

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "challenge_enube/internal/models"
    "challenge_enube/internal/usecase"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrEmailInUse = errors.New("email already in use")

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
