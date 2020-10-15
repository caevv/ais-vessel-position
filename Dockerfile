FROM golang:1.15-alpine as builder

RUN apk --no-cache --update add ca-certificates make

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -o artifacts/svc .

#### Minimal Image #####

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/artifacts/svc svc

EXPOSE 8080

CMD ["./svc"]
