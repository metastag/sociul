package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"math/big"
	"time"

	"sociul-auth-service/internal/models"
	"sociul-auth-service/internal/repository"
	"sociul-auth-service/internal/sentinel"

	"github.com/mailgun/mailgun-go/v5"
)

type OtpService struct {
	repo  *repository.AuthRepository
	cache *repository.Cache
	mg    *mailgun.Client
}

func NewOtpService(repo *repository.AuthRepository, cache *repository.Cache, mg *mailgun.Client) *OtpService {
	return &OtpService{repo: repo, cache: cache, mg: mg}
}

// Create a random 6-digit otp
func createOtp() (string, error) {
	const digits = "0123456789"
	otp := make([]byte, 6)
	for i := range otp {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp[i] = digits[n.Int64()]
	}
	return string(otp), nil
}

// Create a random 32-digit reset token
func createResetToken() (string, error) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Send email to recipient with 10-seconds timeout
func (s *OtpService) sendEmail(recipient, otp string) error {
	sender := "noreply@sociul.com"
	subject := "Sociul Verification Code"
	body := "Your OTP is " + otp

	message := mailgun.NewMessage("", sender, subject, body, recipient)

	emailCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, err := s.mg.Send(emailCtx, message)
	if err != nil {
		return err
	}

	log.Printf("Email sent: %s\n", resp.ID)
	return nil
}

func (s *OtpService) SendEmailOtp(ctx context.Context, request models.SendOtpRequest) error {
	// Create a 6-digit otp
	otp, err := createOtp()
	if err != nil {
		log.Println("Error in SendEmailOtp() createOtp() - ", err)
		return sentinel.ErrInternal
	}

	// Store in cache
	key := "otp:" + request.Email + ":verify"
	err = s.cache.Store(ctx, key, otp, 5*time.Minute)
	if err != nil {
		log.Println("Error in SendEmailOtp() cache Store() - ", err)
		return sentinel.ErrInternal
	}

	// Send otp email to user
	err = s.sendEmail(request.Email, otp)
	if err != nil {
		log.Println("Error in SendEmailOtp() sendEmail() - ", err)
		return sentinel.ErrInternal
	}

	return nil
}

func (s *OtpService) VerifyEmailOtp(ctx context.Context, request models.VerifyOtpRequest) error {
	// Fetch otp from cache
	key := "otp:" + request.Email + ":verify"
	otp, err := s.cache.Fetch(ctx, key)
	if err == sentinel.ErrCacheMiss {
		return err
	}
	if err != nil {
		log.Println("Error in VerifyEmailOtp() cache Fetch() - ", err)
		return sentinel.ErrInternal
	}

	// Match request otp against server otp
	if otp != request.Otp {
		return sentinel.ErrWrongOtp
	}

	// Delete otp from cache
	err = s.cache.Delete(ctx, key)
	if err != nil {
		log.Println("Error in VerifyEmailOtp() cache Delete() - ", err)
		return sentinel.ErrInternal
	}

	// Mark as verified in db
	err = s.repo.MarkVerified(ctx, request.Email)
	if err != nil {
		log.Print("Error in VerifyEmailOtp() MarkVerified() - ", err)
		return sentinel.ErrInternal
	}

	return nil
}

func (s *OtpService) SendResetOtp(ctx context.Context, request models.SendOtpRequest) error {
	// Create a 6-digit otp
	otp, err := createOtp()
	if err != nil {
		log.Println("Error in SendResetOtp() createOtp() - ", err)
		return sentinel.ErrInternal
	}

	// Store in cache
	key := "otp:" + request.Email + ":reset"
	err = s.cache.Store(ctx, key, otp, 5*time.Minute)
	if err != nil {
		log.Println("Error in SendResetOtp() cache Store() - ", err)
		return sentinel.ErrInternal
	}

	// Send otp email to user
	err = s.sendEmail(request.Email, otp)
	if err != nil {
		log.Println("Error in SendResetOtp() sendEmail() - ", err)
		return sentinel.ErrInternal
	}

	return nil
}

func (s *OtpService) VerifyResetOtp(ctx context.Context, request models.VerifyOtpRequest) (string, error) {
	// Fetch otp from cache
	key := "otp:" + request.Email + ":verify"
	otp, err := s.cache.Fetch(ctx, key)
	if err == sentinel.ErrCacheMiss {
		return "", err
	}
	if err != nil {
		log.Println("Error in VerifyResetOtp() cache Fetch() - ", err)
		return "", sentinel.ErrInternal
	}

	// Match request otp against server otp
	if otp != request.Otp {
		return "", sentinel.ErrWrongOtp
	}

	// Delete otp from cache
	err = s.cache.Delete(ctx, request.Email)
	if err != nil {
		log.Println("Error in VerifyResetOtp() cache Delete() - ", err)
		return "", sentinel.ErrInternal
	}

	// Create a reset token
	token, err := createResetToken()
	if err != nil {
		log.Println("Error in VerifyResetOtp() createResetToken() - ", err)
		return "", sentinel.ErrInternal
	}

	// Store token in cache and return it
	key = "reset-token:" + token
	err = s.cache.Store(ctx, key, request.Email, 5*time.Minute)
	if err != nil {
		log.Println("Error in VerifyResetOtp() cache Store() - ", err)
		return "", sentinel.ErrInternal
	}

	return token, nil
}
