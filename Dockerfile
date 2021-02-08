FROM golang:1.15 AS build

WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o main

FROM scratch

WORKDIR /app
COPY --from=build /build/main .
CMD ["./main", "-api", "https://api.warehouse.tld", "-input", "http://localhost/file.csv"]
