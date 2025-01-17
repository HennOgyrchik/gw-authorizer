FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o gw-authorizer ./cmd/

FROM scratch
WORKDIR /app
COPY --from=builder ./app/gw-authorizer .
COPY --from=builder ./app/config.env .
COPY --from=builder ./app/internal/storages/migrations ./migrations
EXPOSE 9090
CMD ["./gw-authorizer"]

#docker build -t gw-authorizer .
#docker run --name gw-authorizer -d gw-authorizer