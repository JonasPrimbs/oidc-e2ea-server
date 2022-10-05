FROM golang:1.19 AS build
WORKDIR /go/src
COPY go ./go
COPY main.go .
COPY go.mod .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o ridt .


FROM scratch AS runtime
COPY --from=build /go/src/ridt ./

ENV ALG="ES256"
ENV DEFAULT_TOKEN_PERIOD=3600
ENV MAX_TOKEN_PERIOD=2592000
ENV PORT=8080

EXPOSE ${PORT}/tcp
ENTRYPOINT ["./ridt"]
