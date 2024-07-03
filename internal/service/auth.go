package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/subbbbbaru/go_final_project/internal/repository"
)

const (
	salt     = "jkjklmklmlkmlkmlmk"
	signKey  = "slcsdcls1!vkmdfk+=/$"
	tokenTTL = 8 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	PwdHash string
}

type AuthService struct {
	repo repository.Auth
}

func NewAuthService(repo repository.Auth) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) GenerateToken(password string) (string, error) {
	pwdFromRepo, err := s.repo.GetPassword()
	if err != nil {
		return "", err
	}
	if pwdFromRepo != password {
		return "", errors.New("password not valid")
	}

	pwdHash := s.genPasswordHash(password)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		pwdHash,
	})
	return token.SignedString([]byte(signKey))
}

func (s *AuthService) ValideToken(accessToken string) (bool, error) {

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("invalid sign method")
		}
		return []byte(signKey), nil
	})
	if err != nil {
		return false, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return false, errors.New("token claims are not of type")
	}

	pwdFromRepo, err := s.repo.GetPassword()
	if err != nil {
		return false, err
	}

	if claims.PwdHash != s.genPasswordHash(pwdFromRepo) {
		return false, errors.New("token not valid")
	}
	return true, nil
}

func (s *AuthService) genPasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
