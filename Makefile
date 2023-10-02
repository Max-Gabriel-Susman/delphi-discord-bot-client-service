
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

# docker build --tag brometheus/delphi-discord-bot-client-service:v0.1.0 .

run: 
	docker run \
		-p 8082:8082 \
		-e API_ADDRESS=0.0.0.0:8082 \
		-e INFERENTIAL_DB_USER=usr \
		-e INFERENTIAL_DB_PASSWORD=identity \
		-e INFERENTIAL_DB_HOST=127.0.0.1 \
		-e INFERENTIAL_DB_NAME=identity \
		-e INFERENTIAL_DB_PORT=3306 \
		-e ENABLE_MIGRATE=true \
		-e BOT_TOKEN=MTEzNzU1MTkxMzYwMTIxMjQ5OA.GXv16P.YRelee4HmWywqcol9Rd_qG9KczGdHpKNr2KZvI \
		brometheus/delphi-discord-bot-client-service:v0.1.2


push: 
	docker push brometheus/delphi-model-service:tagname


update:
	docker build --tag brometheus/delphi-discord-bot-client-service:v0.1.2 .
	docker push brometheus/delphi-discord-bot-client-service:v0.1.2

# docker push brometheus/delphi-discord-bot-client-service:v0.1.0

# grpcurl -plaintext -v localhost:50051 list Greeter

# grpcurl -plaintext -d '{"name": "prometheus"}' localhost:50051 Greeter/SayHello

# grpcurl -plaintext -d '{"name": "prometheus"}' localhost:50052 Greeter/SayHello