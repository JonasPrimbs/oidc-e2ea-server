# Generate build container
FROM golang:1.19 AS build

# Create working directory for source files
WORKDIR /go/src

# Copy source files into container
COPY go ./go
COPY main.go .
COPY go.mod .

# Set compilation environment variables:
# Enable C-Go
ENV CGO_ENABLED=1
# Set target OS to linux
ENV GOOS=linux

# Download dependencies
RUN go get -d -v ./...

# Compile application to single binary file 'iat'
RUN go build -a -ldflags '-linkmode external -extldflags "-static"' -o /go/src/iat


# Generate runtime container
FROM scratch AS runtime

# Create working directory for binary
WORKDIR /

# Copy compiled binary from build container
COPY --from=build /go/src/iat /iat

# Set default configuration parameters
ENV ALG="ES256"
ENV DEFAULT_TOKEN_PERIOD=3600
ENV MAX_TOKEN_PERIOD=2592000
ENV PORT=8080
ENV DB_SQLITE_FILE="/config/db.sqlite"

# Expose the configured TCP port
EXPOSE ${PORT}/tcp

# Define the binary as the entrypoint
ENTRYPOINT ["/iat"]
