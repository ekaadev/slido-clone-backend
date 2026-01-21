# Slido Clone Backend

## Running Application

### Run Web Server

```bash
go run cmd/web/main.go
```

## Database Migrations

Before running the migrations, ensure you have `golang migrate` installed.

```bash

### Create New Migration

If you need to create a new migration, use the following command:

```bash
migrate create -ext sql -dir db/migrations create_table_xxx
```

### Running Migration

If you need to run the migrations, use the following command and create the database `slido_clone` in your MySQL server
beforehand:

```bash
migrate -database "mysql://username:yourpassword@tcp(localhost:3306)/slido_clone" -path db/migrations up
```

## Redis

You need to have a Redis server running locally on the default port `6379`. 