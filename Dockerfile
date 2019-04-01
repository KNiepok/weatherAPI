FROM golang:1.12 as builder

WORKDIR /go/src/github.com/kniepok/weatherAPI

COPY . .

ENV GO111MODULE=on
RUN go mod vendor

# Build the Go app
WORKDIR /go/src/github.com/kniepok/weatherAPI/cmd/main
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/weatherapi .

FROM frolvlad/alpine-glibc

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/bin/weatherapi .

EXPOSE 8080
RUN chmod +x ./weatherapi

CMD ["./weatherapi", "api"]