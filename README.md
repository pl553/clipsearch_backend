# Dependencies
https://github.com/pgvector/pgvector  
https://github.com/openai/CLIP
# Building
```bash
go mod tidy
go build
```
# Testing
```bash
go test ./...
```

# Database
## Initial setup
```bash
sudo -i -u postgres  
createdb clipsearch  
psql clipsearch
```

```
clipsearch=# CREATE EXTENSION vector;
clipsearch=# CREATE USER clipsearch WITH PASSWORD 'weakpassword'
\q
```

```bash 
exit
export POSTGRESQL_URL=postgres://clipsearch:weakpassword@localhost:5432/clipsearch?sslmode=disable
```
## Creating/running migrations on the cli
https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
### Creating a new migration
```bash
migrate create -ext sql -dir db/migrations -seq create_users_table
```
### Running migrations
```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```
