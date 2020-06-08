FROM golang:1.14-alpine AS build

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Copy project files
WORKDIR /go/src
COPY go.mod .
COPY go.sum .

# Fetch dependencies
RUN go mod download
COPY . .

# Build project
ENV CGO_ENABLED=0
RUN go build -o open-dam-api ./cmd/open-dam-api

FROM scratch AS runtime

# Copy dependent files
COPY --from=build /go/src/open-dam-api ./
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8080/tcp

ENTRYPOINT ["./open-dam-api"]
