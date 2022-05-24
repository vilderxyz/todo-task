#name of built binary file 
TODO_BIN = todoApp

# removes old binary and compiles new for linux
build:
	rm -f ./deploy/${TODO_BIN}
	env GOOS=linux CGO_ENABLED=0 go build -o ./deploy/${TODO_BIN} .

# builds binary exec, removes existing containers and sets new ones 
up: build
	docker-compose down
	docker-compose up --build -d

stop: 
	docker-compose stop

start: build
	docker-compose up --build

# stops and removes containers
down:
	docker-compose down

# generate mockModel struct implementing Model interface
mock:
	mockgen -package mock -destination mock/db.go github.com/vilderxyz/todos/db DB

# run temporary database for testing purposes
db:
	docker run --name mock -p 8888:5432 -e POSTGRES_PASSWORD=mock -e POSTGRES_USER=mock -e POSTGRES_DB=mock -d postgres:14.2

# runs tests and remove database container after
test: db
	@echo "Waiting a bit for database to initialize"
	sleep 1
	@echo "Testing database queries..."
	go test -cover  github.com/vilderxyz/todos/db
	@echo "Testing api..."
	go test -cover  github.com/vilderxyz/todos/api
	@echo "Removing temporary database..."
	docker rm -f mock

.PHONY: db test mock build up down
