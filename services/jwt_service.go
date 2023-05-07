package services

import (
	"fmt"
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
	Sum uint `json:"sum"`
	jwt.StandardClaims
}

func (s *jwtService) GenerateToken(id uint) (string, error) {
	TOKEN_HOURS_LIFESPAN := "1"
	token_lifespan, err := strconv.Atoi(TOKEN_HOURS_LIFESPAN)
	if err != nil {
		return "", err
	}

	claim := &Claim{
		id,
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

func (s *jwtService) ValidateToken(token string) bool {
	_, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("Invalid token: %v", token)
		}
		return []byte(s.secretKey), nil
	})

	return err == nil
}
