# Building
```bash
go mod tidy
go build
```
# Database
## Initial setup
```bash
sudo -i -u postgres  
psql  
```

```
postgres=# create user clipsearch with password 'weakpassword'
\q
```

```bash
createdb clipsearch  
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
