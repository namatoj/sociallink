
# Build the web app
FROM golang:1.22.0-alpine AS builder
ARG PB_VERSION=0.21.3

WORKDIR $GOPATH/sociallink

COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o dist/main cmd/web/main.go

# Runner
FROM gcr.io/distroless/static
COPY --from=builder go/sociallink/dist/main .
EXPOSE 8080
CMD ["./main", "serve", "--http=0.0.0.0:8080"]
