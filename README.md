# Database
## Initial setup
`sudo -i -u postgres`  
`psql`  
`postgres=# create user clipsearch with password 'weakpassword'`    
`\q`   
`createdb clipsearch`  
`exit`  
`export POSTGRESQL_URL=postgres://clipsearch:weakpassword@localhost:5432/clipsearch?sslmode=disable`  
## Creating/running migrations on the cli
https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
### Creating a new migration
`migrate create -ext sql -dir db/migrations -seq create_users_table`
### Running migrations
`migrate -database ${POSTGRESQL_URL} -path db/migrations up` 
