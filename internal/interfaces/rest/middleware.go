package rest

import (
	"context"
	"github.com/018bf/companies/internal/entity"
	"github.com/018bf/companies/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//go:generate mockgen -source=middleware.go -package=rest -destination=middleware_mock.go
type authInterceptor interface {
	ValidateToken(ctx context.Context, token *entity.Token) error
}

type AuthMiddleware struct {
	authService authInterceptor
}

func NewAuthMiddleware(authService authInterceptor) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		header := c.GetHeader("Authorization")
		var token *entity.Token
		if len(header) > 7 {
			header = header[7:]
			token = entity.NewToken(header)
			if err := m.authService.ValidateToken(ctx, token); err != nil {
				decodeError(c, err)
				return
			}
		}
		ctx = context.WithValue(ctx, TokenContextKey, token)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func RequestMiddleware(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), log.RequestIDKey, uuid.New().String())
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}
