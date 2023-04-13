package integration

import (
	"context"
	"flag"
	"fmt"
	"github.com/018bf/companies/internal/configs"
	"github.com/018bf/companies/internal/containers"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/interfaces/postgres"
	postgres_repositories "github.com/018bf/companies/internal/repositories/postgres"
	"github.com/jaswdr/faker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		exitVal := m.Run()
		os.Exit(exitVal)
		return
	}
	ctx := context.Background()
	configPath := os.Getenv("COMPANIES_CONFIG_PATH")
	var err error
	config, err = configs.ParseConfig(configPath)
	if err != nil {
		panic(err)
	}
	migrate := containers.NewMigrateContainer(configPath)
	if err := migrate.Start(ctx); err != nil {
		panic(err)
	}
	app := containers.NewGRPCContainer(configPath)
	go func() {
		if err := app.Start(ctx); err != nil {
			panic(err)
		}
	}()
	db, err := postgres.NewDatabase(config)
	if err != nil {
		panic(err)
	}
	database = db
	conn, err = grpc.Dial(config.BindAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	companyRepository = postgres_repositories.NewCompanyRepository(database, nil)
	company = &models.Company{
		ID:                "",
		UpdatedAt:         faker.New().Time().Time(time.Now()),
		CreatedAt:         faker.New().Time().Time(time.Now()),
		Name:              faker.New().Lorem().Text(15),
		Description:       faker.New().Lorem().Sentence(15),
		AmountOfEmployees: faker.New().IntBetween(1, 1000),
		Registered:        faker.New().Bool(),
		Type:              models.CompanyType(faker.New().Int8Between(1, 4)),
	}
	err = companyRepository.Create(ctx, company)
	if err != nil {
		panic(err)
	}
	exitVal := m.Run()
	tables := []string{
		"public.companies",
	}
	_, err = database.ExecContext(
		ctx,
		fmt.Sprintf("TRUNCATE TABLE %s CASCADE", strings.Join(tables, ", ")),
	)
	if err != nil {
		panic(err)
	}
	os.Exit(exitVal)
}
