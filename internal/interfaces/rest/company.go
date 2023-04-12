package rest

import (
	"fmt"
	"net/http"

	"github.com/018bf/companies/internal/domain/interceptors"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/pkg/log"
	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	companyInterceptor interceptors.CompanyInterceptor
	logger             log.Logger
}

func NewCompanyHandler(
	companyInterceptor interceptors.CompanyInterceptor,
	logger log.Logger,
) *CompanyHandler {
	return &CompanyHandler{companyInterceptor: companyInterceptor, logger: logger}
}

func (h *CompanyHandler) Register(router *gin.RouterGroup) {
	group := router.Group("/companies")
	group.POST("/", h.Create)
	group.GET("/", h.List)
	group.GET("/:id", h.Get)
	group.PATCH("/:id", h.Update)
	group.DELETE("/:id", h.Delete)
}

// Create        godoc
// @Summary      Store a new Company
// @Description  Takes a Company JSON and store in DB. Return saved JSON.
// @Tags         Company
// @Produce      json
// @Param        Company  body   models.CompanyCreate  true  "Company JSON"
// @Success      201   {object}  models.Company
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /companies/ [post]
func (h *CompanyHandler) Create(ctx *gin.Context) {
	token := ctx.Request.Context().Value(TokenContextKey).(*models.Token)
	create := &models.CompanyCreate{}
	_ = ctx.Bind(create)
	company, err := h.companyInterceptor.Create(ctx.Request.Context(), create, token)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, company)
}

// List          godoc
// @Summary      List Company array
// @Description  Responds with the list of all Company as JSON.
// @Tags         Company
// @Produce      json
// @Param        filter  query   models.CompanyFilter false "Company filter"
// @Success      200  {array}  models.Company
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /companies [get]
func (h *CompanyHandler) List(ctx *gin.Context) {
	token := ctx.Request.Context().Value(TokenContextKey).(*models.Token)
	filter := &models.CompanyFilter{}
	_ = ctx.Bind(filter)
	listCompanies, count, err := h.companyInterceptor.List(
		ctx.Request.Context(),
		filter,
		token,
	)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.Header("count", fmt.Sprint(count))
	ctx.JSON(http.StatusOK, listCompanies)
}

// Get           godoc
// @Summary      Get single Company by UUID
// @Description  Returns the Company whose UUID value matches the UUID.
// @Tags         Company
// @Produce      json
// @Param        uuid  path      string  true  "search Company by UUID"
// @Success      200  {object}  models.Company
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /companies/{uuid} [get]
func (h *CompanyHandler) Get(ctx *gin.Context) {
	token := ctx.Request.Context().Value(TokenContextKey).(*models.Token)
	company, err := h.companyInterceptor.Get(
		ctx.Request.Context(),
		models.UUID(ctx.Param("id")),
		token,
	)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, company)
}

// Update        godoc
// @Summary      Update Company by UUID
// @Description  Returns the updated Company.
// @Tags         Company
// @Produce      json
// @Param        uuid  path      string  true  "update Company by UUID"
// @Param        Company  body   models.CompanyUpdate  true  "Company JSON"
// @Success      201  {object}  models.Company
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /companies/{uuid} [PATCH]
func (h *CompanyHandler) Update(ctx *gin.Context) {
	token := ctx.Request.Context().Value(TokenContextKey).(*models.Token)
	update := &models.CompanyUpdate{}
	_ = ctx.Bind(update)
	update.ID = models.UUID(ctx.Param("id"))
	company, err := h.companyInterceptor.Update(ctx.Request.Context(), update, token)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, company)
}

// Delete        godoc
// @Summary      Delete single Company by UUID
// @Description  Delete the Company whose UUID value matches the UUID.
// @Tags         Company
// @Param        uuid  path      string  true  "delete Company by UUID"
// @Success      204
// @Failure      400   {object}  errs.Error
// @Failure      401   {object}  errs.Error
// @Failure      403   {object}  errs.Error
// @Failure      404   {object}  errs.Error
// @Failure      405   {object}  errs.Error
// @Failure      500   {object}  errs.Error
// @Failure      503   {object}  errs.Error
// @Router       /companies/{uuid} [delete]
func (h *CompanyHandler) Delete(ctx *gin.Context) {
	token := ctx.Request.Context().Value(TokenContextKey).(*models.Token)
	err := h.companyInterceptor.Delete(
		ctx.Request.Context(),
		models.UUID(ctx.Param("id")),
		token,
	)
	if err != nil {
		decodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
