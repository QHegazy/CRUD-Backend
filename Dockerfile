FROM golang:1.24.2-bookworm AS builder

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./cmd/api/main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app


COPY --from=builder /app/app .

USER nonroot:nonroot


EXPOSE 8080

CMD ["./app"]