FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o vanieth .

FROM scratch

COPY --from=builder /app/vanieth /vanieth

ENTRYPOINT ["/vanieth"]
