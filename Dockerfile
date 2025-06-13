# Dockerfile

# --- Build Stage ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Salin file dependensi terlebih dahulu untuk cache
COPY go.mod go.sum ./
RUN go mod download

# Salin SEMUA file dan folder proyek (termasuk cmd/ dan internal/)
COPY . .

# Kompilasi aplikasi dengan menunjuk ke direktori main package Anda
# Go akan secara otomatis menemukan semua paket lain di dalam 'internal/'
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/api

# --- Final Stage ---
FROM alpine:3.18

# Install paket tzdata untuk data timezone
RUN apk add --no-cache tzdata

WORKDIR /app

# Salin hanya binary yang sudah jadi dari 'builder' stage
COPY --from=builder /app/main .

EXPOSE 3000

# Perintah untuk menjalankan aplikasi
CMD ["/app/main"]
