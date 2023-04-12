package rest

import (
	"context"
	"github.com/018bf/companies/internal/domain/interceptors"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	authService interceptors.AuthInterceptor
}

func NewAuthMiddleware(authService interceptors.AuthInterceptor) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		header := c.GetHeader("Authorization")
		var token *models.Token
		if len(header) > 7 {
			header = header[7:]
			token = models.NewToken(header)
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
