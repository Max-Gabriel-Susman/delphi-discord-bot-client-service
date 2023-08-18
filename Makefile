
mod: 
	go mod tidy
	go mod vendor 

local-start:
	go run cmd/delphi-discord-bot-client-service/main.go