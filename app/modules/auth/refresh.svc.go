package auth

import (
	"context"

	"github.com/google/uuid"
)

type RefreshTokenService struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenServiceResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) RefreshTokenService(ctx context.Context, req *RefreshTokenService) (*RefreshTokenServiceResponse, error) {
	claims, err := s.parseToken(req.RefreshToken, "refresh")
	if err != nil {
		return nil, ErrInvalidToken
	}
	memberID, err := uuid.Parse(claims.Sub)
	if err != nil {
		return nil, ErrInvalidToken
	}

	accessToken, _, err := s.generateToken(memberID, claims.Email, claims.Role, claims.IsAdmin, "access", accessTokenTTL)
	if err != nil {
		return nil, err
	}
	refreshToken, _, err := s.generateToken(memberID, claims.Email, claims.Role, claims.IsAdmin, "refresh", refreshTokenTTL)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenServiceResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	}, nil
}
