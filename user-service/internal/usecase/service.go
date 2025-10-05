package usecase

import (
	"errors"
	"time"

	domain "github.com/adopabianko/commerce/user-service/internal/domain/user"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo      domain.Repository
	jwtSecret string
}

func New(repo domain.Repository, jwtSecret string) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s *Service) Register(email, password, name string) (*domain.User, error) {
	existing, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &domain.User{Email: email, Password: string(hash), Name: name}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Login(email, password string) (string, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	// create token
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	signed, err := t.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *Service) ValidateToken(tokenStr string) (uint, error) {
	tok, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := tok.Claims.(jwt.MapClaims); ok && tok.Valid {
		subf, ok := claims["sub"].(float64)
		if !ok {
			return 0, errors.New("invalid sub claim")
		}
		return uint(subf), nil
	}
	return 0, errors.New("invalid token")
}
