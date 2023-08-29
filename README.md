# CLIP search backend
[OpenAI CLIP](https://github.com/openai/CLIP) based image search backend. This program maintains a set of images that you can search through with text prompts. Functionality is exposed through a HTTP API.  
### [API Documentation](https://pl553.github.io/clipsearch_api_redoc/)
# Setup
Install [CLIP](https://github.com/openai/CLIP), [pgvector](https://github.com/pgvector/pgvector), [migrate](https://github.com/golang-migrate/migrate), libzmq and libsodium.  

Create the database to be used:
```
sudo -i -u postgres
createdb clipsearch  
psql clipsearch
```
Install the vector extension into the database and create a user:
```
CREATE EXTENSION vector;
CREATE USER clipsearch WITH PASSWORD 'weakpassword'
```
The connection url:
```bash 
export POSTGRESQL_URL=postgres://clipsearch:weakpassword@localhost:5432/clipsearch?sslmode=disable
```
Run the database migrations:
```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```
Build the program:
```bash
go mod tidy
go build
```
Run one of the CLIP daemons and let it download the model (about 900 MB). The model will be downloaded into the models/ relative directory.
```bash
cd clip_daemons
python image_embedding_daemon.py
```

# Usage
Launch the clip daemons.
```bash
cd clip_daemons
python image_embedding_daemon.py &
python text_embedding_daemon.py &
```
Launch the program:
```bash
export POSTGRESQL_URL=postgres://clipsearch:weakpassword@localhost:5432/clipsearch?sslmode=disable
./clipsearch
```
The server is now listening on port 3000.  
See https://github.com/pl553/clipsearch/ on how this is integrated with a frontend.
# Environment variables
| Variable | Meaning | Default |
| --- | --- | --- |
| PORT | The port that the http server will listen on | 3000
| POSTGRESQL_URL | Database connection url | - |
| ZMQ_IMAGE_PORT | The port that the image embedding daemon is expected to be on. The program will attempt to connect to tcp://localhost:${ZMQ_IMAGE_PORT} over zmq | 5554 |
| ZMQ_TEXT_PORT | The port that the text embedding daemon is expected to be on. | 5553

# Testing
```bash
go test ./...
```
# Creating/running migrations on the cli
https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
### Creating a new migration
```bash
migrate create -ext sql -dir db/migrations -seq create_users_table
```
### Running migrations
```bash
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```
