package service

import (
	"context"

	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/log"
)

//go:generate mockgen -source=auth.go -package=usecases -destination=auth_mock.go

type authRepository interface {
	Validate(ctx context.Context, token *entity.Token) error
	GetSubject(ctx context.Context, token *entity.Token) (string, error)
	HasPermission(
		ctx context.Context,
		permission entity.PermissionID,
		token *entity.Token,
	) error
	HasObjectPermission(
		ctx context.Context,
		permission entity.PermissionID,
		token *entity.Token,
		obj any,
	) error
}

type AuthService struct {
	authRepository authRepository
	logger         log.Logger
}

func NewAuthService(
	authRepository authRepository,
	logger log.Logger,
) *AuthService {
	return &AuthService{
		authRepository: authRepository,
		logger:         logger,
	}
}

func (u AuthService) ValidateToken(ctx context.Context, token *entity.Token) error {
	if err := u.authRepository.Validate(ctx, token); err != nil {
		return err
	}
	return nil
}

func (u AuthService) HasPermission(
	ctx context.Context,
	token *entity.Token,
	permission entity.PermissionID,
) error {
	if err := u.authRepository.HasPermission(ctx, permission, token); err != nil {
		return err
	}
	return nil
}

func (u AuthService) HasObjectPermission(
	ctx context.Context,
	token *entity.Token,
	permission entity.PermissionID,
	object any,
) error {
	if err := u.authRepository.HasObjectPermission(ctx, permission, token, object); err != nil {
		return err
	}
	return nil
}
