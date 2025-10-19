# Football Team API

API **RESTful** yang dibangun menggunakan **Go (Golang)**, **Gin**, dan **GORM** untuk mengelola tim sepak bola, pemain, pertandingan, dan hasil pertandingan.  
Proyek ini mendukung manajemen admin, laporan pertandingan, dan dilengkapi seeder untuk data awal.

---

## Kebutuhan

- Go >= 1.21
- MySQL
- Postman (opsional, untuk pengujian API)

---

## Setup

1. Clone repository:

```bash
git clone
cd footballteam
```

2. Install dependency:
```bash
go mod tidy
```

3. Copy .env.example menjadi .env dan ubah sesuai kebutuhan:
DB_USER=root
DB_PASSWORD=
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=footballteam
APP_PORT=8080

## Menjalankan Proyek
```bash
go run main.go
```

## Dokumentasi API Postman
https://documenter.getpostman.com/view/9770363/2sB3QQJo4R
