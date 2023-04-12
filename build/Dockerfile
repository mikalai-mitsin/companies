FROM golang:1.20 as builder

WORKDIR /go/src/companies
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags "-X github.com/018bf/companies.Version=$(git describe --tags --abbrev=0)" -v -o ./dist/companies ./cmd/companies

FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/companies/dist/companies /app/companies
ENTRYPOINT ["/app/companies"]
