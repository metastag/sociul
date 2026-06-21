package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"sociul-auth-service/internal/models"
	"sociul-auth-service/internal/repository"
	"sociul-auth-service/internal/sentinel"
	"sociul-auth-service/internal/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *repository.AuthRepository
	cache  *repository.Cache
	secret string
}

func NewAuthService(repo *repository.AuthRepository, cache *repository.Cache, secret string) *AuthService {
	return &AuthService{repo: repo, cache: cache, secret: secret}
}

// Generate refresh token
func (s *AuthService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Generate jwt access token with 15min ttl and user id as claim
func (s *AuthService) generateAccessToken(id uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":  id,
		"role": "user", // For now everyone is a user
		"exp":  time.Now().Add(time.Minute * 15).Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Verify jwt token and return user id claim
func (s *AuthService) ValidateJwt(tokenString string) (uuid.UUID, error) {
	// Verify jwt token string
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, sentinel.ErrInvalidToken
	}

	// Extract user id claim and write to context
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, sentinel.ErrInvalidToken
	}

	id, err := claims.GetSubject()
	if err != nil {
		return uuid.Nil, sentinel.ErrInvalidToken
	}

	userId, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, sentinel.ErrInvalidToken

	}
	return userId, nil
}

func (s *AuthService) SignUp(ctx context.Context, request models.SignUpRequest) (string, string, error) {
	// Normalize Input
	request = utils.SignUpNormalize(request)

	// Validation check
	err := utils.SignUpValidate(request)
	if err != nil {
		return "", "", err
	}

	// Check if username exists
	exists, err := s.repo.UserExists(ctx, request.Username)
	if err != nil {
		log.Println("Error in SignUp() UserExists() call - ", err)
		return "", "", sentinel.ErrInternal
	}
	if exists {
		return "", "", sentinel.ErrUsernameAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 12)
	if err != nil {
		log.Println("Error in SignUp() while hashing password - ", err)
		return "", "", sentinel.ErrInternal
	}
	request.Password = string(hashedPassword)

	// Generate uuid
	id := uuid.New()

	// Create user dto
	var user models.User
	user.Id = id
	user.Username = request.Username
	user.Password = request.Password
	user.Email = request.Email

	// Send to repository to create user
	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Println("Error in SignUp() CreateUser() call - ", err)
		return "", "", sentinel.ErrInternal
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		log.Println("Error in SignUp() generateRefreshToken() - ", err)
		return "", "", sentinel.ErrInternal
	}

	// Store refresh token in cache
	key := "refresh-token:" + refreshToken
	err = s.cache.Store(ctx, key, id.String(), 240*time.Hour)

	// Generate access token (short ttl)
	accessToken, err := s.generateAccessToken(id)
	if err != nil {
		log.Println("Error in SignUp() generateAccessToken() - ", err)
		return "", "", sentinel.ErrInternal
	}

	return refreshToken, accessToken, nil
}

func (s *AuthService) Login(ctx context.Context, request models.LoginRequest) (string, string, error) {
	// Normalize Input
	request = utils.LoginNormalize(request)

	// Validation check
	err := utils.LoginValidate(request)
	if err != nil {
		return "", "", err
	}

	// Retrieve password and uuid of username
	id, password, err := s.repo.GetUserCreds(ctx, request.Username)
	if err == sentinel.ErrUserNotFound {
		return "", "", err
	} else if err != nil {
		log.Println("Error in Login() GetUserCreds() call - ", err)
		return "", "", sentinel.ErrInternal
	}

	// Compare hashed passwords
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(request.Password))
	if err != nil {
		return "", "", sentinel.ErrWrongPassword
	}

	// Generate refresh token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		log.Println("Error in Login() generateRefreshToken() - ", err)
		return "", "", sentinel.ErrInternal
	}

	// Store refresh token -> uuid in cache
	key := "refresh-token:" + refreshToken
	err = s.cache.Store(ctx, key, id.String(), 240*time.Hour)

	// Generate access token (short ttl)
	accessToken, err := s.generateAccessToken(id)
	if err != nil {
		log.Println("Error in Login() generateAccessToken() - ", err)
		return "", "", sentinel.ErrInternal
	}

	return refreshToken, accessToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, request models.RefreshRequest) (string, error) {
	// Fetch value from cache
	key := "refresh-token:" + request.RefreshToken

	value, err := s.cache.Fetch(ctx, key)
	if err == sentinel.ErrCacheMiss {
		return "", err
	} else if err != nil {
		log.Println("Error in Refresh() Fetch() call - ", err)
		return "", sentinel.ErrInternal
	}

	// Parse string value as uuid
	id, err := uuid.Parse(value)
	if err != nil {
		return "", sentinel.ErrInvalidToken
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(id)
	if err != nil {
		log.Println("Error in Refresh() generateAccessToken() - ", err)
		return "", sentinel.ErrInternal
	}

	return accessToken, nil
}

func (s *AuthService) Logout(ctx context.Context, request models.RefreshRequest) error {
	key := "refresh-token:" + request.RefreshToken

	// Delete refresh token from cache
	err := s.cache.Delete(ctx, key)
	if err != nil {
		log.Println("Error in Logout() Delete() call - ", err)
		return sentinel.ErrInternal
	}
	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, request models.ResetPasswordRequest) error {
	key := "reset-token:" + request.ResetToken
	email, err := s.cache.Fetch(ctx, key)
	if err == sentinel.ErrCacheMiss {
		return err
	} else if err != nil {
		log.Println("Error in ResetPassword() Fetch() call - ", err)
		return sentinel.ErrInternal
	}

	// Delete reset token from cache
	err = s.cache.Delete(ctx, key)
	if err != nil {
		log.Println("Error in ResetPassword() Delete() call - ", err)
		return sentinel.ErrInternal
	}

	// Normalize and validate password
	password := request.Password
	password = strings.TrimSpace(password)
	err = utils.ValidatePassword(password)
	if err != nil {
		return err
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Println("Error in ResetPassword() while hashing password - ", err)
		return sentinel.ErrInternal
	}

	err = s.repo.UpdatePassword(ctx, email, string(hashedPassword))
	if err != nil {
		log.Println("Error in ResetPassword() UpdatePassword() call - ", err)
		return sentinel.ErrInternal
	}
	return nil
}

func (s *AuthService) DeleteUser(ctx context.Context) error {
	// Extract user id from context
	id, ok := utils.FromContext(ctx)
	if !ok {
		return sentinel.ErrUnauthorized
	}

	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		log.Println("Error in DeleteUser() - ", err)
		return sentinel.ErrInternal
	}
	return nil

}
