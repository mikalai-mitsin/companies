package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/018bf/companies/internal/configs"
	"github.com/018bf/companies/internal/errs"
	"github.com/018bf/companies/pkg/log"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ctxKey string

const (
	TokenContextKey ctxKey = "token"
)

type Server struct {
	router *gin.Engine
	config *configs.Config
	logger log.Logger
}

// NewServer        godoc
// @title           companies
// @version         0.1.0
// @description     TBD
// @host      127.0.0.1:8000
// @BasePath  /api/v1
// @schemes https http
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
//
// @security ApiKeyAuth
func NewServer(
	logger log.Logger,
	config *configs.Config,
	authMiddleware *AuthMiddleware,
	companyHandler *CompanyHandler,
) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(authMiddleware.Middleware())
	router.Use(cors.Default())
	router.Use(RequestMiddleware)
	router.Use(Logger(logger))
	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})
	apiV1 := router.Group("api").Group("v1")
	companyHandler.Register(apiV1)
	return &Server{
		router: router,
		config: config,
		logger: logger,
	}
}

func (s *Server) Start(_ context.Context) error {
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) Stop(_ context.Context) error {
	return nil
}

func decodeError(ctx *gin.Context, err error) {
	var domainError *errs.Error
	if errors.As(err, &domainError) {
		switch domainError.Code {
		case errs.ErrorCodeOK:
			ctx.JSON(http.StatusOK, err)
		case errs.ErrorCodeCanceled:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeUnknown:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeInvalidArgument:
			ctx.JSON(http.StatusBadRequest, err)
		case errs.ErrorCodeDeadlineExceeded:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeNotFound:
			ctx.JSON(http.StatusNotFound, err)
		case errs.ErrorCodeAlreadyExists:
			ctx.JSON(http.StatusBadRequest, err)
		case errs.ErrorCodePermissionDenied:
			ctx.JSON(http.StatusForbidden, err)
		case errs.ErrorCodeResourceExhausted:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeFailedPrecondition:
			ctx.JSON(http.StatusBadRequest, err)
		case errs.ErrorCodeAborted:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeOutOfRange:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeUnimplemented:
			ctx.JSON(http.StatusMethodNotAllowed, err)
		case errs.ErrorCodeInternal:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeUnavailable:
			ctx.JSON(http.StatusServiceUnavailable, err)
		case errs.ErrorCodeDataLoss:
			ctx.JSON(http.StatusInternalServerError, err)
		case errs.ErrorCodeUnauthenticated:
			ctx.JSON(http.StatusUnauthorized, err)
		default:
			ctx.JSON(http.StatusInternalServerError, err)
		}
		return
	}
	ctx.String(http.StatusInternalServerError, err.Error())
}

func Logger(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		end := time.Now().UTC()
		latency := end.Sub(start)
		fields := []log.Field{
			log.Int("status", c.Writer.Status()),
			log.String("method", c.Request.Method),
			log.String("path", path),
			log.String("query", query),
			log.String("ip", c.ClientIP()),
			log.String("user-agent", c.Request.UserAgent()),
			log.Duration("latency", latency),
			log.String("time", end.String()),
			log.Context(c.Request.Context()),
		}
		if len(c.Errors) > 0 {
			// Append error field if this is an erroneous request.
			for _, e := range c.Errors.Errors() {
				logger.Error(e, fields...)
			}
		} else {
			logger.Info(path, fields...)
		}
	}
}
