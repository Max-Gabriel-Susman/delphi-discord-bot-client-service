
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
	go run cmd/delphi-discord-bot-client-service/main.go