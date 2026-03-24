# Generate build container
FROM golang:1.26 AS build

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

# Compile application to single binary file 'ict'
RUN go build -a -ldflags '-linkmode external -extldflags "-static"' -o /go/src/oidc2middleware


# Generate runtime container
FROM scratch AS runtime

# Create working directory for binary
WORKDIR /

# Copy compiled binary from build container
COPY --from=build /go/src/oidc2middleware /bin/oidc2middleware

# Set default configuration parameters
ENV KEY_FILE="/run/secrets/private_key"
ENV KEY_ID="1"
ENV ISSUER="https://op.localhost/realms/oidc2"
ENV ALG="ES256"
ENV TOKEN_PERIOD="300"
ENV PORT="8080"

# Expose the configured TCP port
EXPOSE ${PORT}/tcp

# Define the binary as the entrypoint
ENTRYPOINT ["/bin/oidc2middleware"]
