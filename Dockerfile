FROM golang:1.19 AS build
WORKDIR /go/src
COPY go ./go
COPY main.go .
COPY go.mod .

ENV CGO_ENABLED=1
ENV GOOS=linux
RUN go get -d -v ./...

RUN go build -a -ldflags '-linkmode external -extldflags "-static"' -o /go/src/ridt


FROM scratch AS runtime
COPY --from=build /go/src/ridt /ridt

ENV ALG="ES256"
ENV DEFAULT_TOKEN_PERIOD=3600
ENV MAX_TOKEN_PERIOD=2592000
ENV PORT=8080
ENV DB_SQLITE_FILE="/config/db.sqlite"

WORKDIR /

EXPOSE ${PORT}/tcp
ENTRYPOINT ["/ridt"]
