package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/018bf/companies/internal/configs"
	"github.com/018bf/companies/internal/domain/errs"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/repositories"
	"github.com/018bf/companies/pkg/clock"
	"github.com/018bf/companies/pkg/log"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const accessAudience = "access"

type AuthRepository struct {
	accessTTL  time.Duration
	refreshTTL time.Duration
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	clock      clock.Clock
	logger     log.Logger
}

func NewAuthRepository(
	config *configs.Config,
	clock clock.Clock,
	logger log.Logger,
) repositories.AuthRepository {
	private, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(config.Auth.PrivateKey))
	if err != nil {
		panic(err)
	}
	public, err := jwt.ParseRSAPublicKeyFromPEM([]byte(config.Auth.PublicKey))
	if err != nil {
		panic(err)
	}
	return &AuthRepository{
		accessTTL:  time.Duration(config.Auth.AccessTTL) * time.Second,
		refreshTTL: time.Duration(config.Auth.RefreshTTL) * time.Second,
		publicKey:  public,
		privateKey: private,
		clock:      clock,
		logger:     logger,
	}
}

func (r *AuthRepository) Validate(_ context.Context, token *models.Token) error {
	jwtToken, err := r.parse(token)
	if err != nil {
		return err
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	if !claims.VerifyAudience(accessAudience, true) {
		return errs.NewBadToken()
	}
	return nil
}

func (r *AuthRepository) GetSubject(ctx context.Context, token *models.Token) (string, error) {
	jwtToken, err := r.parse(token)
	if err != nil {
		e := errs.NewError(errs.ErrorCodeUnauthenticated, "Invalid token.")
		return "", e
	}
	claims := jwtToken.Claims.(jwt.MapClaims)
	return fmt.Sprint(claims["sub"]), nil
}

func (r *AuthRepository) parse(token *models.Token) (*jwt.Token, error) {
	if token == nil {
		return nil, errs.NewBadToken()
	}
	jwtToken, err := jwt.Parse(token.String(), r.keyFunc)
	if err != nil {
		e := errs.NewBadToken()
		return nil, e
	}
	return jwtToken, nil
}
func (r *AuthRepository) keyFunc(_ *jwt.Token) (interface{}, error) {
	return r.publicKey, nil
}

type objectPermissionChecker func(model any, token *jwt.Token) error
type permissionChecker func(token *jwt.Token) error

var hasObjectPermission = map[models.PermissionID][]objectPermissionChecker{
	models.PermissionIDCompanyList:   {objectAnybody},
	models.PermissionIDCompanyDetail: {objectAnybody},
	models.PermissionIDCompanyCreate: {objectUser},
	models.PermissionIDCompanyUpdate: {objectUser},
	models.PermissionIDCompanyDelete: {objectUser},
}

var hasPermission = map[models.PermissionID][]permissionChecker{
	models.PermissionIDCompanyList:   {anybody},
	models.PermissionIDCompanyDetail: {anybody},
	models.PermissionIDCompanyCreate: {user},
	models.PermissionIDCompanyUpdate: {user},
	models.PermissionIDCompanyDelete: {user},
}

func (r *AuthRepository) HasPermission(
	_ context.Context,
	permissionID models.PermissionID,
	token *models.Token,
) error {
	t, _ := r.parse(token)
	checkers := hasPermission[permissionID]
	for _, checker := range checkers {
		if err := checker(t); err == nil {
			return nil
		}
	}
	return errs.NewPermissionDenied()
}

func (r *AuthRepository) HasObjectPermission(
	_ context.Context,
	permission models.PermissionID,
	token *models.Token,
	obj any,
) error {
	t, _ := r.parse(token)
	checkers := hasObjectPermission[permission]
	for _, checker := range checkers {
		if err := checker(obj, t); err == nil {
			return nil
		}
	}
	return errs.NewPermissionDenied()
}

// nolint: unused
func objectAdmin(_ any, token *jwt.Token) error {
	if token == nil {
		return errs.NewPermissionDenied()
	}
	claims := token.Claims.(jwt.MapClaims)
	isAdminClaims, contains := claims["admin"]
	if isAdmin, ok := isAdminClaims.(bool); contains && ok && isAdmin {
		return nil
	}
	return errs.NewPermissionDenied()
}

func objectUser(_ any, token *jwt.Token) error {
	if token == nil {
		return errs.NewPermissionDenied()
	}
	return nil
}

// nolint: unused
func objectNobody(_ any, _ *jwt.Token) error {
	return errs.NewPermissionDenied()
}

func objectAnybody(_ any, _ *jwt.Token) error {
	return nil
}

// nolint: unused
func admin(token *jwt.Token) error {
	if token == nil {
		return errs.NewPermissionDenied()
	}
	claims := token.Claims.(jwt.MapClaims)
	isAdminClaims, contains := claims["admin"]
	if isAdmin, ok := isAdminClaims.(bool); contains && ok && isAdmin {
		return nil
	}
	return errs.NewPermissionDenied()
}

func user(token *jwt.Token) error {
	if token == nil {
		return errs.NewPermissionDenied()
	}
	return nil
}

// nolint: unused
func nobody(_ *jwt.Token) error {
	return errs.NewPermissionDenied()
}

func anybody(_ *jwt.Token) error {
	return nil
}
