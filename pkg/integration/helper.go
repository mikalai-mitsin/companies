package integration

import (
	"context"
	"errors"
	"fmt"
	"github.com/018bf/companies/internal/configs"
	"github.com/018bf/companies/internal/domain/models"
	"github.com/018bf/companies/internal/domain/repositories"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var accessToken = models.NewToken("eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNjRiOGIwZC00ZDg1LTRiMjMtOWQ3Ny0wNTRmMDY5MGU0OTgiLCJhdWQiOlsiYWNjZXNzIl0sImV4cCI6MTc0NzY1MjQzOCwibmJmIjoxNjYxMjUyNDM4LCJpYXQiOjE2NjEyNTI0MzgsImp0aSI6ImViY2M0MDUwLTU3YzMtNGVlMy1hNjMzLWY1NzgyOTc0MjRjYSJ9.N2amwawNdnpgVNcZq4LOwtcDK88USnilTPeH79Dvv10oHU2QW4hC4t68n7LcbPRWyX-ZwvhcpWAq3xaTkGNP0vmvuPGUJrPcwbxmqSdXLrrg7U-xH2tyXtdSZeZwtYgSp3D9haQXTm74S2fqNDvhhSx28Pp_3uSsXMYzxS4R2mo")
var (
	company           *models.Company
	database          *sqlx.DB
	config            *configs.Config
	conn              *grpc.ClientConn
	companyRepository repositories.CompanyRepository
)

func getAuthContext(ctx context.Context, token *models.Token) context.Context {
	md := metadata.New(map[string]string{"authorization": fmt.Sprintf("Bearer %s", token)})
	return metadata.NewOutgoingContext(ctx, md)
}

func statusEqual(x, y *status.Status) bool {
	x, _ = x.WithDetails(nil)
	y, _ = y.WithDetails(nil)
	return errors.Is(x.Err(), y.Err())
}
