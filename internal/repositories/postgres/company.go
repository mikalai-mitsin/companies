package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/018bf/companies/internal/domain/errs"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/repositories"
	"github.com/018bf/companies/pkg/log"
	"github.com/018bf/companies/pkg/postgresql"
	"github.com/018bf/companies/pkg/utils"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type CompanyRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewCompanyRepository(database *sqlx.DB, logger log.Logger) repositories.CompanyRepository {
	return &CompanyRepository{database: database, logger: logger}
}
func (r *CompanyRepository) Create(ctx context.Context, company *models.Company) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewCompanyDTOFromModel(company)
	q := sq.Insert("public.companies").
		Columns(
			"updated_at",
			"created_at",
			"name",
			"description",
			"amount_of_employees",
			"registered",
			"type",
		).
		Values(
			dto.UpdatedAt,
			dto.CreatedAt,
			dto.Name,
			dto.Description,
			dto.AmountOfEmployees,
			dto.Registered,
			dto.Type,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	company.ID = models.UUID(dto.ID)
	return nil
}
func (r *CompanyRepository) Get(ctx context.Context, id models.UUID) (*models.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &CompanyDTO{}
	q := sq.Select(
		"companies.id",
		"companies.updated_at",
		"companies.created_at",
		"companies.name",
		"companies.description",
		"companies.amount_of_employees",
		"companies.registered",
		"companies.type",
	).
		From("public.companies").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("company_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}

func (r *CompanyRepository) List(
	ctx context.Context,
	filter *models.CompanyFilter,
) ([]*models.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto CompanyListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = utils.Pointer(pageSize)
	}
	q := sq.Select(
		"companies.id",
		"companies.updated_at",
		"companies.created_at",
		"companies.name",
		"companies.description",
		"companies.amount_of_employees",
		"companies.registered",
		"companies.type",
	).
		From("public.companies").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(
			postgresql.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"name", "description"},
			},
		)
	}
	if len(filter.IDs) > 0 {
		q = q.Where(sq.Eq{"id": filter.IDs})
	}
	if len(filter.Types) > 0 {
		q = q.Where(sq.Eq{"type": filter.Types})
	}
	if filter.Registered != nil {
		q = q.Where(sq.Eq{"registered": *filter.Registered})
	}
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	q = q.Limit(*filter.PageSize)
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &dto, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return dto.ToModels(), nil
}

func (r *CompanyRepository) Count(
	ctx context.Context,
	filter *models.CompanyFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.companies")
	if filter.Search != nil {
		q = q.Where(
			postgresql.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"name", "description"},
			},
		)
	}
	if len(filter.IDs) > 0 {
		q = q.Where(sq.Eq{"id": filter.IDs})
	}
	if len(filter.Types) > 0 {
		q = q.Where(sq.Eq{"type": filter.Types})
	}
	if filter.Registered != nil {
		q = q.Where(sq.Eq{"registered": *filter.Registered})
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result := r.database.QueryRowxContext(ctx, query, args...)
	if err := result.Err(); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	var count uint64
	if err := result.Scan(&count); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	return count, nil
}
func (r *CompanyRepository) Update(ctx context.Context, company *models.Company) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewCompanyDTOFromModel(company)
	q := sq.Update("public.companies").Where(sq.Eq{"id": company.ID}).
		Set("updated_at", dto.UpdatedAt).
		Set("name", dto.Name).
		Set("description", dto.Description).
		Set("amount_of_employees", dto.AmountOfEmployees).
		Set("registered", dto.Registered).
		Set("type", dto.Type)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("company_id", fmt.Sprint(company.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("company_id", fmt.Sprint(company.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("company_id", fmt.Sprint(company.ID))
		return e
	}
	return nil
}
func (r *CompanyRepository) Delete(ctx context.Context, id models.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.companies").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("company_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("company_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().WithParam("company_id", fmt.Sprint(id))
		return e
	}
	return nil
}

type CompanyDTO struct {
	ID                string    `db:"id,omitempty"`
	UpdatedAt         time.Time `db:"updated_at,omitempty"`
	CreatedAt         time.Time `db:"created_at,omitempty"`
	Name              string    `db:"name"`
	Description       string    `db:"description"`
	AmountOfEmployees int       `db:"amount_of_employees"`
	Registered        bool      `db:"registered"`
	Type              uint8     `db:"type"`
}
type CompanyListDTO []*CompanyDTO

func (list CompanyListDTO) ToModels() []*models.Company {
	listCompanies := make([]*models.Company, len(list))
	for i := range list {
		listCompanies[i] = list[i].ToModel()
	}
	return listCompanies
}
func NewCompanyDTOFromModel(company *models.Company) *CompanyDTO {
	dto := &CompanyDTO{
		ID:                string(company.ID),
		UpdatedAt:         company.UpdatedAt,
		CreatedAt:         company.CreatedAt,
		Name:              company.Name,
		Description:       company.Description,
		AmountOfEmployees: company.AmountOfEmployees,
		Registered:        company.Registered,
		Type:              uint8(company.Type),
	}
	return dto
}
func (dto *CompanyDTO) ToModel() *models.Company {
	model := &models.Company{
		ID:                models.UUID(dto.ID),
		UpdatedAt:         dto.UpdatedAt,
		CreatedAt:         dto.CreatedAt,
		Name:              dto.Name,
		Description:       dto.Description,
		AmountOfEmployees: dto.AmountOfEmployees,
		Registered:        dto.Registered,
		Type:              models.CompanyType(dto.Type),
	}
	return model
}
