# Slido Clone Backend

Backend service untuk aplikasi Slido Clone dengan fitur Q&A, Polling, Live Chat, dan Leaderboard.

## Prerequisites

- Go 1.21+
- MySQL 8.0+
- Redis 7.0+

## Getting Started

### 1. Install Dependencies

#### Install Golang Migrate

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

**macOS:**
```bash
brew install golang-migrate
```

**Windows:**
```bash
scoop install migrate
```

#### Install Redis

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install redis-server
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

**macOS:**
```bash
brew install redis
brew services start redis
```

**Windows:**
Download dari [Redis for Windows](https://github.com/microsoftarchive/redis/releases) atau gunakan WSL.

### 2. Setup Database

Access MySQL shell:
```bash
mysql -u root -p
```

Create user and database:
```sql
CREATE USER 'slido_user'@'localhost' IDENTIFIED BY 'password';
CREATE DATABASE slido_clone;
GRANT ALL PRIVILEGES ON slido_clone.* TO 'slido_user'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

Jalankan migration (pastikan nama database sesuai dengan `slido_clone`):
```bash
migrate -database "mysql://slido_user:password@tcp(localhost:3306)/slido_clone" -path db/migrations up
```

> **Note:** Ganti `password` dengan password yang aman untuk production.

### 3. Install Go Dependencies

```bash
go mod tidy
go build ./...
```

### 4. Setup Configuration

Buat file config (bebas, bisa `.env` atau `config.json` tergantung implementasi `config.NewViper`). Pastikan database name adalah `slido_clone`.

Contoh jika menggunakan environment variables:

```bash
export DB_NAME=slido_clone
export DB_USER=slido_user
export DB_PASSWORD=password
export DB_HOST=localhost
export DB_PORT=3306
```

### 5. Run Application

```bash
go run cmd/web/main.go
```

Server akan berjalan di `http://localhost:3000` (atau sesuai konfigurasi).

### Check Redis Connection

```bash
redis-cli ping
# Expected output: PONG
```