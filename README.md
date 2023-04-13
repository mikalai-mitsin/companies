# companies

### Docs
- [`Changelog`](docs/CHANGELOG.md)

### Run with `docker-compose`
`docker-compose -f deployments/docker-compose.yaml up`

### Usage:
`companies [global options] command [command options] [arguments...]`

### Commands:
- `migrate`  Run migrations
- `grpc`     Run gRPC server
- `rest`     Run REST server
- `help`, `h`  Shows a list of commands or help for one command

### Global options:
- `--config FILE`, `-c FILE`  Load configuration from FILE [$COMPANIES_CONFIG_PATH]
- `--help`, `-h`              show help

## Unit Testing
1. Run `task unit`

## Integration Testing
1. Create tests database
2. Create kafka topic
3. Fill `configs/test.toml`
4. Run `task test`

## To improve
- JWT permissions checking сan be simplified if we don't plan to check permissions in the token
- Remove deprecated list method (left for example)
- Add DTOs to REST module

## Directories

### `/cmd`

Main applications for this project.

### `/internal`

Private application and library code.

## Service Application Directories

### `/api`

OpenAPI/Swagger specs, JSON schema files, protocol definition files.

## Common Application Directories

### `/deployments`

Configurations for deploy the project to servers.

## Other Directories

### `/docs`

Design and user documents (in addition to your godoc generated documentation).
