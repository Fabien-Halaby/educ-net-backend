package usecase

import (
	"educnet/internal/auth"
	"educnet/internal/domain"
	"educnet/internal/handler/dto"
	"educnet/internal/repository"
	"errors"
)

type AuthUseCase interface {
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
}

type authUseCase struct {
	userRepo   repository.UserRepository
	jwtService *auth.JWTService
}

func NewAuthUseCase(userRepo repository.UserRepository, jwtService *auth.JWTService) AuthUseCase {
	return &authUseCase{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *authUseCase) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	//! 1. Find user by email
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	//! 2. Verify password
	if !user.VerifyPassword(req.Password) {
		return nil, domain.ErrInvalidCredentials
	}

	//! 3. Check if user is approved (only approved users can login)
	if !user.IsApproved() {
		return nil, errors.New("your account is pending approval")
	}

	//! 4. Generate tokens
	accessToken, err := uc.jwtService.GenerateAccessToken(
		user.ID,
		user.Email,
		user.Role,
		user.SchoolID,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	//! 5. Return response
	return &dto.LoginResponse{
		User: dto.UserInfo{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.GetFullName(),
			Role:      user.Role,
			Status:    user.Status,
			SchoolID:  user.SchoolID,
			AvatarURL: user.AvatarURL,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    24 * 60 * 60, //! 24 heures en secondes
	}, nil
}
