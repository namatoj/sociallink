# sociallink

## Requirements

 - [Go](https://go.dev/) version 1.22+

## Local development

Run `./scripts/init-repo.sh` to set up commit hooks. Install the dependencies with `go mod tidy` Install [air](https://github.com/cosmtrek/air) to enable live reloading.

## Web app

With `air` installed start the app by running `air`. Otherwise start the application with `go run cmd/web/main.go serve`.

## CLI

Run cli commands with `go run cmd/cli/main.go <command>`. e.g `go run cmd/cli/main.go admin create <email> <password>`

## Conventional commits

This project follows conventional commits. Read more about it [here](https://www.conventionalcommits.org/en/v1.0.0/).

