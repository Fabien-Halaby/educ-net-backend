package auth

import (
	"errors"
	"time"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	SchoolID int    `json:"school_id"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewJWTService(secret string, accessTTL, refreshTTL int) *JWTService {
	return &JWTService{
		secretKey:       secret,
		accessTokenTTL:  time.Duration(accessTTL) * time.Hour,
		refreshTokenTTL: time.Duration(refreshTTL) * 24 * time.Hour,
	}
}

//! GenerateAccessToken génère un access token JWT
func (s *JWTService) GenerateAccessToken(userID int, email, role string, schoolID int) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		Role:     role,
		SchoolID: schoolID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

//! GenerateRefreshToken génère un refresh token JWT
func (s *JWTService) GenerateRefreshToken(userID int, email string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

//! ValidateToken valide et parse un token JWT
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	log.Printf("[JWT] Validating token with secret: %s...", s.secretKey[:10]) // DEBUG
	
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Vérifier la méthode de signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("[JWT] Invalid signing method: %v", token.Method) // DEBUG
			return nil, ErrInvalidToken
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		log.Printf("[JWT] Parse error: %v", err) // DEBUG
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		log.Printf("[JWT] Invalid claims or token not valid") // DEBUG
		return nil, ErrInvalidToken
	}

	// Vérifier l'expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		log.Printf("[JWT] Token expired at: %v", claims.ExpiresAt.Time) // DEBUG
		return nil, ErrExpiredToken
	}

	log.Printf("[JWT] Token valid! User ID: %d, Email: %s", claims.UserID, claims.Email) // DEBUG
	return claims, nil
}

