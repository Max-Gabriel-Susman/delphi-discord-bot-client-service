
nuke-deps:
	rm go.mod go.sum
	rm -r vendor 

mod-init:
	go mod init github.com/Max-Gabriel-Susman/delphi-discord-bot-client-service
	go mod tidy
	go mod vendor

mod: 
	go mod tidy
	go mod vendor 

local-start:
	go run main.go

build:
	docker build --tag delphi-discord-bot-client-service .

run: 
	docker run \
		-e API_ADDRESS=0.0.0.0:8082 \
		-e INFERENTIAL_DB_USER=usr \
		-e INFERENTIAL_DB_PASSWORD=identity \
		-e INFERENTIAL_DB_HOST=127.0.0.1 \
		-e INFERENTIAL_DB_NAME=identity \
		-e INFERENTIAL_DB_PORT=3306 \
		-e ENABLE_MIGRATE=true \
		-e BOT_TOKEN=MTEzNzU1MTkxMzYwMTIxMjQ5OA.G0zj61.MQLAarnh_XR64lPg8RI5fZur1JKTi2JIaqaiX4 \
		delphi-discord-bot-client-service
