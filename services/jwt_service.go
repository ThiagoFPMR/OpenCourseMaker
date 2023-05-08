package services

import (
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() *jwtService {
	return &jwtService{
		secretKey: "secret-key",
		issuer:    "OpenCourseMaker",
	}
}

type Claim struct {
	ID    uint   `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func (s *jwtService) GenerateToken(id uint, nome string, email string) (string, error) {
	TOKEN_HOURS_LIFESPAN := "1"
	token_lifespan, err := strconv.Atoi(TOKEN_HOURS_LIFESPAN)
	if err != nil {
		return "", err
	}

	claim := &Claim{
		id,
		nome,
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix(),
			Issuer:    s.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}
