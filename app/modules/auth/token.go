package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type tokenClaims struct {
	Sub          string `json:"sub"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	IsAdmin      bool   `json:"is_admin"`
	ActorSub     string `json:"actor_sub,omitempty"`
	ActorIsAdmin bool   `json:"actor_is_admin,omitempty"`
	ActingAs     bool   `json:"acting_as,omitempty"`
	Typ          string `json:"typ"`
	Iat          int64  `json:"iat"`
	Exp          int64  `json:"exp"`
}

type tokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

const accessTokenTTL = 24 * time.Hour
const refreshTokenTTL = 7 * 24 * time.Hour

func (s *Service) generateToken(memberID uuid.UUID, email string, role string, isAdmin bool, actorSub string, actorIsAdmin bool, actingAs bool, tokenType string, ttl time.Duration) (string, int64, error) {
	if s.secret == "" {
		return "", 0, errors.New("missing token secret")
	}

	now := time.Now().UTC()
	exp := now.Add(ttl).Unix()

	header := tokenHeader{Alg: "HS256", Typ: "JWT"}
	claims := tokenClaims{
		Sub:          memberID.String(),
		Email:        email,
		Role:         role,
		IsAdmin:      isAdmin,
		ActorSub:     actorSub,
		ActorIsAdmin: actorIsAdmin,
		ActingAs:     actingAs,
		Typ:          tokenType,
		Iat:          now.Unix(),
		Exp:          exp,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", 0, err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", 0, err
	}

	enc := base64.RawURLEncoding
	headerB64 := enc.EncodeToString(headerJSON)
	claimsB64 := enc.EncodeToString(claimsJSON)
	signingInput := headerB64 + "." + claimsB64

	mac := hmac.New(sha256.New, []byte(s.secret))
	_, _ = mac.Write([]byte(signingInput))
	sig := enc.EncodeToString(mac.Sum(nil))

	return signingInput + "." + sig, exp, nil
}

func (s *Service) parseToken(token string, expectedType string) (*tokenClaims, error) {
	if s.secret == "" {
		return nil, ErrInvalidToken
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	enc := base64.RawURLEncoding
	signingInput := parts[0] + "." + parts[1]

	mac := hmac.New(sha256.New, []byte(s.secret))
	_, _ = mac.Write([]byte(signingInput))
	expected := mac.Sum(nil)

	provided, err := enc.DecodeString(parts[2])
	if err != nil || !hmac.Equal(expected, provided) {
		return nil, ErrInvalidToken
	}

	claimsJSON, err := enc.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims tokenClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.Exp <= time.Now().UTC().Unix() {
		return nil, ErrInvalidToken
	}
	if expectedType != "" && claims.Typ != expectedType {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}
